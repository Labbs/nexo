package auth

import (
	"fmt"
	"time"

	"github.com/labbs/nexo/application/auth/dto"
	s "github.com/labbs/nexo/application/session/dto"
	u "github.com/labbs/nexo/application/user/dto"
	"github.com/labbs/nexo/infrastructure/helpers/tokenutil"
	"golang.org/x/crypto/bcrypt"
)

func (c *AuthApp) Authenticate(input dto.AuthenticateInput) (*dto.AuthenticateOutput, error) {
	logger := c.Logger.With().Str("component", "application.auth.authenticate").Logger()

	resp, err := c.UserApp.GetByEmail(u.GetByEmailInput{Email: input.Email})
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if !resp.User.Active {
		logger.Warn().Str("email", input.Email).Msg("attempt to authenticate inactive user")
		return nil, fmt.Errorf("user is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(input.Password))
	if err != nil {
		logger.Warn().Str("email", input.Email).Msg("invalid password attempt")
		return nil, fmt.Errorf("invalid credentials")
	}

	sessionResult, err := c.SessionApp.Create(s.CreateSessionInput{
		UserId:    resp.User.Id,
		UserAgent: input.Context.Get("User-Agent"),
		IpAddress: input.Context.IP(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(c.Config.Session.ExpirationMinutes)),
	})
	if err != nil {
		logger.Error().Err(err).Str("user_id", resp.User.Id).Msg("failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	accessToken, err := tokenutil.CreateAccessToken(resp.User.Id, sessionResult.SessionId, c.Config)
	if err != nil {
		logger.Error().Err(err).Str("user_id", resp.User.Id).Str("session_id", sessionResult.SessionId).Msg("failed to create access token")
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	return &dto.AuthenticateOutput{
		Token: accessToken,
	}, nil
}
