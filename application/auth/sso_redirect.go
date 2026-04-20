package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/labbs/nexo/application/auth/dto"
	"golang.org/x/oauth2"
)

func (c *AuthApplication) SSORedirect() (*dto.SSORedirectOutput, error) {
	if !c.Config.SSO.Enabled {
		return nil, fmt.Errorf("SSO is not enabled")
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	// state = base64(nonce) + "." + base64(HMAC(nonce))
	nonce := base64.RawURLEncoding.EncodeToString(raw)
	mac := hmac.New(sha256.New, []byte(c.Config.Session.SecretKey))
	mac.Write([]byte(nonce))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	state := nonce + "." + sig

	oauthCfg := c.buildOAuthConfig()
	url := oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOnline)

	return &dto.SSORedirectOutput{
		URL:   url,
		State: state,
	}, nil
}
