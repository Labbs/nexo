package session

import (
	"fmt"
	"time"

	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/session/dto"
	databaseDto "github.com/labbs/nexo/application/database/dto"
	drawingDto "github.com/labbs/nexo/application/drawing/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	userDto "github.com/labbs/nexo/application/user/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/infrastructure/helpers/tokenutil"
)

func (c *SessionApplication) ValidateToken(input dto.ValidateTokenInput) (*dto.ValidateTokenOutput, error) {
	logger := c.Logger.With().Str("component", "application.session.validate_token").Logger()

	sessionId, err := tokenutil.GetSessionIdFromToken(input.Token, c.Config)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get session id from token")
		return nil, apperrors.ErrInvalidToken
	}

	session, err := c.SessionPers.GetById(sessionId)
	if err != nil {
		logger.Error().Err(err).Str("session_id", sessionId).Msg("failed to get session by id")
		return nil, apperrors.ErrInvalidToken
	}

	if session.ExpiresAt.Before(time.Now()) {
		logger.Warn().Str("session_id", sessionId).Msg("session has expired")
		return nil, apperrors.ErrSessionExpired
	}

	// Fetch user to populate roles
	userResult, err := c.UserApplication.GetByUserId(userDto.GetByUserIdInput{UserId: session.UserId})
	if err != nil {
		logger.Error().Err(err).Str("user_id", session.UserId).Msg("failed to get user for role population")
		return nil, apperrors.ErrInvalidToken
	}

	ctx := &fiberoapi.AuthContext{
		UserID: session.UserId,
		Roles:  []string{string(userResult.User.Role)},
		Claims: map[string]any{
			"session_id": session.Id,
		},
	}

	return &dto.ValidateTokenOutput{AuthContext: ctx}, nil
}

func (c *SessionApplication) HasRole(input dto.HasRoleInput) bool {
	for _, r := range input.Context.Roles {
		if r == input.Role {
			return true
		}
	}
	return false
}

func (c *SessionApplication) HasScope(input dto.HasScopeInput) bool {
	for _, s := range input.Context.Scopes {
		if s == input.Scope {
			return true
		}
	}
	return false
}

func (c *SessionApplication) CanAccessResource(input dto.CanAccessResourceInput) (bool, error) {
	authCtx := input.Context

	// Admin bypasses all resource access checks
	for _, r := range authCtx.Roles {
		if r == string(domain.RoleAdmin) {
			return true, nil
		}
	}

	requiredRole := actionToRequiredRole(input.Action)

	switch input.ResourceType {
	case "space":
		return c.canAccessSpace(authCtx.UserID, input.ResourceID, requiredRole)
	case "database":
		return c.canAccessDatabase(authCtx.UserID, input.ResourceID, requiredRole)
	case "drawing":
		return c.canAccessDrawing(authCtx.UserID, input.ResourceID, requiredRole)
	default:
		return false, nil
	}
}

func actionToRequiredRole(action string) string {
	switch action {
	case "read":
		return "viewer"
	case "create", "write":
		return "editor"
	case "delete":
		return "owner"
	default:
		return "viewer"
	}
}

func (c *SessionApplication) canAccessSpace(userID, spaceID, requiredRole string) (bool, error) {
	result, err := c.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: spaceID})
	if err != nil {
		return false, err
	}
	return result.Space.HasPermission(userID, requiredRole), nil
}

func (c *SessionApplication) canAccessDatabase(userID, databaseID, requiredRole string) (bool, error) {
	result, err := c.DatabaseApplication.GetDatabaseById(databaseDto.GetDatabaseByIdInput{DatabaseId: databaseID})
	if err != nil {
		return false, err
	}
	return c.canAccessSpace(userID, result.Database.SpaceId, requiredRole)
}

func (c *SessionApplication) canAccessDrawing(userID, drawingID, requiredRole string) (bool, error) {
	result, err := c.DrawingApplication.GetDrawingById(drawingDto.GetDrawingByIdInput{DrawingId: drawingID})
	if err != nil {
		return false, err
	}
	return c.canAccessSpace(userID, result.Drawing.SpaceId, requiredRole)
}

func (c *SessionApplication) GetUserPermissions(input dto.GetUserPermissionsInput) (*dto.GetUserPermissionsOutput, error) {
	logger := c.Logger.With().Str("component", "application.session.get_user_permissions").Logger()

	logger.Warn().Msg("not implemented")

	return nil, fmt.Errorf("not implemented")
}
