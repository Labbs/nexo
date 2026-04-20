package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labbs/nexo/application/auth/dto"
	d "github.com/labbs/nexo/application/document/dto"
	s "github.com/labbs/nexo/application/session/dto"
	spdto "github.com/labbs/nexo/application/space/dto"
	u "github.com/labbs/nexo/application/user/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/tokenutil"
	"golang.org/x/oauth2"
)

type oidcUserInfo struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
}

func (c *AuthApplication) SSOCallback(input dto.SSOCallbackInput) (*dto.SSOCallbackOutput, error) {
	logger := c.Logger.With().Str("component", "application.auth.sso_callback").Logger()

	if !c.Config.SSO.Enabled {
		return nil, fmt.Errorf("SSO is not enabled")
	}

	if err := c.verifyState(input.State); err != nil {
		logger.Warn().Err(err).Msg("invalid SSO state")
		return nil, fmt.Errorf("invalid state parameter")
	}

	oauthCfg := c.buildOAuthConfig()
	token, err := oauthCfg.Exchange(context.Background(), input.Code)
	if err != nil {
		logger.Error().Err(err).Msg("failed to exchange OAuth code")
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	userInfo, err := c.fetchUserInfo(oauthCfg, token)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch userinfo")
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	if userInfo.Sub == "" {
		return nil, fmt.Errorf("provider did not return a user identifier")
	}

	user, err := c.findOrCreateSSOUser(userInfo)
	if err != nil {
		logger.Error().Err(err).Str("sub", userInfo.Sub).Msg("failed to find or create SSO user")
		return nil, fmt.Errorf("failed to resolve user: %w", err)
	}

	sessionResult, err := c.SessionApplication.Create(s.CreateSessionInput{
		UserId:    user.Id,
		UserAgent: input.Context.Get("User-Agent"),
		IpAddress: input.Context.IP(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(c.Config.Session.ExpirationMinutes)),
	})
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.Id).Msg("failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	accessToken, err := tokenutil.CreateAccessToken(user.Id, sessionResult.SessionId, c.Config)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.Id).Msg("failed to create access token")
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	return &dto.SSOCallbackOutput{Token: accessToken}, nil
}

// verifyState validates the HMAC-signed state parameter.
func (c *AuthApplication) verifyState(state string) error {
	parts := strings.SplitN(state, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("malformed state")
	}
	nonce, sig := parts[0], parts[1]
	mac := hmac.New(sha256.New, []byte(c.Config.Session.SecretKey))
	mac.Write([]byte(nonce))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return fmt.Errorf("state signature mismatch")
	}
	return nil
}

// fetchUserInfo calls the provider's userinfo endpoint using the access token.
func (c *AuthApplication) fetchUserInfo(oauthCfg *oauth2.Config, token *oauth2.Token) (*oidcUserInfo, error) {
	// Prefer the endpoint discovered via OIDC; fall back to /userinfo.
	endpoints := []string{strings.TrimRight(c.Config.SSO.IssuerURL, "/") + "/userinfo"}
	if c.oidcUserinfoEndpoint != "" && c.oidcUserinfoEndpoint != endpoints[0] {
		endpoints = append([]string{c.oidcUserinfoEndpoint}, endpoints...)
	}

	client := oauthCfg.Client(context.Background(), token)
	var lastErr error
	for _, url := range endpoints {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("userinfo endpoint returned %d", resp.StatusCode)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var info oidcUserInfo
		if err := json.Unmarshal(body, &info); err != nil {
			return nil, err
		}
		return &info, nil
	}
	return nil, lastErr
}

// findOrCreateSSOUser finds an existing user linked to the SSO provider, or creates a new one.
func (c *AuthApplication) findOrCreateSSOUser(info *oidcUserInfo) (domain.User, error) {
	// 1. Check if provider link already exists
	op, err := c.OAuthProviderPers.FindByProviderAndSubject("oidc", info.Sub)
	if err == nil {
		// Link exists — fetch the user
		resp, err := c.UserApplication.GetByUserId(u.GetByUserIdInput{UserId: op.UserId})
		if err != nil {
			return domain.User{}, err
		}
		return *resp.User, nil
	}

	// 2. Try to link an existing user by email
	var user domain.User
	if info.Email != "" {
		resp, err := c.UserApplication.GetByEmail(u.GetByEmailInput{Email: info.Email})
		if err == nil {
			user = *resp.User
		}
	}

	// 3. No existing user — auto-create one
	if user.Id == "" {
		username := info.PreferredUsername
		if username == "" {
			username = strings.Split(info.Email, "@")[0]
		}
		if username == "" {
			username = "user-" + info.Sub[:8]
		}

		created, err := c.UserApplication.Create(u.CreateUserInput{
			User: domain.User{
				Username: username,
				Email:    info.Email,
				Password: "", // no password for SSO users
				Active:   true,
			},
		})
		if err != nil {
			return domain.User{}, fmt.Errorf("failed to create SSO user: %w", err)
		}
		user = *created.User

		// Create private space + welcome document (mirrors Register use case)
		space, err := c.SpaceApplication.CreatePrivateSpaceForUser(spdto.CreatePrivateSpaceForUserInput{UserId: user.Id})
		if err == nil {
			welcomeContent := []d.Block{{
				ID:   "welcome-1",
				Type: d.BlockTypeParagraph,
				Props: map[string]any{
					"textColor": "default", "backgroundColor": "default", "textAlignment": "left",
				},
				Content: []d.InlineContent{{
					Type: "text", Text: "This is your private space. Start adding your notes and documents here!",
					Styles: map[string]bool{},
				}},
				Children: []d.Block{},
			}}
			_, _ = c.DocumentApplication.CreateDocument(d.CreateDocumentInput{
				Name:    "Welcome to Your Private Space",
				UserId:  user.Id,
				SpaceId: space.Space.Id,
				Content: welcomeContent,
			})
		}
	}

	// 4. Create the provider link
	_, err = c.OAuthProviderPers.Create(domain.OAuthProvider{
		UserId:         user.Id,
		Provider:       "oidc",
		ProviderUserId: info.Sub,
		Email:          info.Email,
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to store SSO provider link: %w", err)
	}

	return user, nil
}
