package auth

import (
	"fmt"

	"github.com/labbs/nexo/application/auth/dto"
	s "github.com/labbs/nexo/application/session/dto"
)

func (c *AuthApp) Logout(input dto.LogoutInput) error {
	logger := c.Logger.With().Str("component", "application.auth.logout").Logger()

	err := c.SessionApp.InvalidateSession(s.InvalidateSessionInput{SessionId: input.SessionId})
	if err != nil {
		logger.Error().Err(err).Str("session_id", input.SessionId).Msg("failed to invalidate session")
		return fmt.Errorf("failed to invalidate session: %w", err)
	}

	return nil
}
