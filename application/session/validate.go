package session

import (
	"fmt"
	"time"

	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/session/dto"
	"github.com/labbs/nexo/infrastructure/helpers/tokenutil"
)

func (c *SessionApplication) ValidateToken(input dto.ValidateTokenInput) (*dto.ValidateTokenOutput, error) {
	logger := c.Logger.With().Str("component", "application.session.validate_token").Logger()

	sessionId, err := tokenutil.GetSessionIdFromToken(input.Token, c.Config)
	if err != nil {
		logger.Error().Err(err).Str("token", input.Token).Msg("failed to get session id from token")
		return nil, apperrors.ErrInvalidToken
	}

	session, err := c.SessionPers.GetById(sessionId)
	if err != nil {
		logger.Error().Err(err).Str("token", input.Token).Msg("failed to get session by token")
		return nil, apperrors.ErrInvalidToken
	}

	if session.ExpiresAt.Before(time.Now()) {
		logger.Warn().Str("token", input.Token).Msg("session has expired")
		return nil, apperrors.ErrSessionExpired
	}

	ctx := &fiberoapi.AuthContext{
		UserID: session.UserId,
		Claims: map[string]any{
			"session_id": session.Id,
		},
	}

	return &dto.ValidateTokenOutput{AuthContext: ctx}, nil
}

func (c *SessionApplication) HasRole(input dto.HasRoleInput) bool {
	logger := c.Logger.With().Str("component", "application.session.has_role").Logger()

	logger.Warn().Msg("not implemented")

	return false
}

func (c *SessionApplication) HasScope(input dto.HasScopeInput) bool {
	logger := c.Logger.With().Str("component", "application.session.has_scope").Logger()

	logger.Warn().Msg("not implemented")

	return false
}

func (c *SessionApplication) CanAccessResource(input dto.CanAccessResourceInput) (bool, error) {
	logger := c.Logger.With().Str("component", "application.session.can_access_resource").Logger()

	logger.Warn().Msg("not implemented")

	return false, fmt.Errorf("not implemented")
}

func (c *SessionApplication) GetUserPermissions(input dto.GetUserPermissionsInput) (*dto.GetUserPermissionsOutput, error) {
	logger := c.Logger.With().Str("component", "application.session.get_user_permissions").Logger()

	logger.Warn().Msg("not implemented")

	return nil, fmt.Errorf("not implemented")
}
