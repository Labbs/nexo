package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type oidcDiscovery struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
}

// discoverOIDC fetches the provider's discovery document and caches the endpoints.
func (c *AuthApplication) discoverOIDC(ctx context.Context) (*oidcDiscovery, error) {
	discoveryURL := strings.TrimRight(c.Config.SSO.IssuerURL, "/") + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OIDC discovery request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OIDC discovery returned %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var d oidcDiscovery
	if err := json.Unmarshal(body, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (c *AuthApplication) buildOAuthConfig() *oauth2.Config {
	authURL := c.Config.SSO.AuthURL
	tokenURL := c.Config.SSO.TokenURL

	// If explicit URLs not provided, attempt discovery at startup.
	// Errors here are non-fatal; the callback will fail gracefully.
	if authURL == "" || tokenURL == "" {
		if disc, err := c.discoverOIDC(context.Background()); err == nil {
			if authURL == "" {
				authURL = disc.AuthorizationEndpoint
			}
			if tokenURL == "" {
				tokenURL = disc.TokenEndpoint
			}
			// Cache the userinfo endpoint for use in fetchUserInfo.
			c.oidcUserinfoEndpoint = disc.UserinfoEndpoint
		} else {
			c.Logger.Warn().Err(err).Msg("OIDC discovery failed — check SSO config")
		}
	}

	scopes := []string{"openid", "email", "profile"}
	scopes = append(scopes, c.Config.SSO.Scopes...)

	return &oauth2.Config{
		ClientID:     c.Config.SSO.ClientID,
		ClientSecret: c.Config.SSO.ClientSecret,
		RedirectURL:  c.Config.SSO.RedirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
}
