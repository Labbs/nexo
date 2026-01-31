package space

import (
	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/interfaces/http/v1/space/dtos"
)

func (ctrl *Controller) CreateSpace(ctx *fiber.Ctx, req dtos.CreateSpaceRequest) (*dtos.CreateSpaceResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.space.create_space").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	// Default to public if not specified
	spaceType := domain.SpaceTypePublic
	if req.Type != nil && *req.Type == string(domain.SpaceTypePrivate) {
		spaceType = domain.SpaceTypePrivate
	}

	result, err := ctrl.SpaceApp.CreateSpace(spaceDto.CreateSpaceInput{
		Name:      req.Name,
		Icon:      req.Icon,
		IconColor: req.IconColor,
		OwnerId:   &authCtx.UserID,
		Type:      spaceType,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create space")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to create space",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.CreateSpaceResponse{
		SpaceId: result.Space.Id,
	}, nil
}

// GetMySpaces handled in user routes (/user/my-spaces)

func (ctrl *Controller) UpdateSpace(ctx *fiber.Ctx, req dtos.UpdateSpaceRequest) (*dtos.UpdateSpaceResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.space.update_space").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.SpaceApp.UpdateSpace(spaceDto.UpdateSpaceInput{
		UserId:    authCtx.UserID,
		SpaceId:   req.SpaceId,
		Name:      req.Name,
		Icon:      req.Icon,
		IconColor: req.IconColor,
	})
	if err != nil {
		// Permission vs not found vs generic
		if err.Error() == "forbidden" {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if err.Error() == "not_found" {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Space not found", Type: "SPACE_NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update space")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update space", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.UpdateSpaceResponse{Space: dtos.Space{
		Id:        result.Space.Id,
		Name:      result.Space.Name,
		Slug:      result.Space.Slug,
		Icon:      result.Space.Icon,
		IconColor: result.Space.IconColor,
		Type:      string(result.Space.Type),
		CreatedAt: result.Space.CreatedAt,
		UpdatedAt: result.Space.UpdatedAt,
	}}, nil
}

func (ctrl *Controller) DeleteSpace(ctx *fiber.Ctx, req dtos.DeleteSpaceRequest) (*dtos.DeleteSpaceResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.space.delete_space").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if err := ctrl.SpaceApp.DeleteSpace(spaceDto.DeleteSpaceInput{
		UserId:  authCtx.UserID,
		SpaceId: req.SpaceId,
	}); err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "conflict_children":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusConflict, Details: "Space has active documents", Type: "CONFLICT"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Space not found", Type: "SPACE_NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to delete space")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete space", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.DeleteSpaceResponse{SpaceId: req.SpaceId}, nil
}

func (ctrl *Controller) ListPermissions(ctx *fiber.Ctx, req dtos.ListSpacePermissionsRequest) (*dtos.ListSpacePermissionsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.space.list_permissions").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.SpaceApp.ListSpacePermissions(spaceDto.ListSpacePermissionsInput{
		UserId:  authCtx.UserID,
		SpaceId: req.SpaceId,
	})
	if err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Space not found", Type: "SPACE_NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to list permissions")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list permissions", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	resp := &dtos.ListSpacePermissionsResponse{Permissions: make([]dtos.SpacePermission, len(result.Permissions))}
	for i, p := range result.Permissions {
		resp.Permissions[i] = dtos.SpacePermission{
			Id:      p.Id,
			UserId:  p.UserId,
			GroupId: p.GroupId,
			Role:    string(p.Role),
		}
		if p.User != nil {
			resp.Permissions[i].Username = p.User.Username
		}
		if p.Group != nil {
			resp.Permissions[i].GroupName = p.Group.Name
		}
	}
	return resp, nil
}

func (ctrl *Controller) UpsertUserPermission(ctx *fiber.Ctx, req dtos.UpsertSpaceUserPermissionRequest) (*dtos.UpsertSpaceUserPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.space.upsert_user_permission").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	var role domain.PermissionRole
	switch req.Role {
	case string(domain.PermissionRoleOwner):
		role = domain.PermissionRoleOwner
	case string(domain.PermissionRoleAdmin):
		role = domain.PermissionRoleAdmin
	case string(domain.PermissionRoleEditor):
		role = domain.PermissionRoleEditor
	default:
		role = domain.PermissionRoleViewer
	}

	if err := ctrl.SpaceApp.UpsertSpaceUserPermission(spaceDto.UpsertSpaceUserPermissionInput{
		RequesterId:  authCtx.UserID,
		SpaceId:      req.SpaceId,
		TargetUserId: req.UserId,
		Role:         role,
	}); err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Space not found", Type: "SPACE_NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to upsert permission")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to upsert permission", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.UpsertSpaceUserPermissionResponse{Message: "permission updated"}, nil
}

func (ctrl *Controller) DeleteUserPermission(ctx *fiber.Ctx, req dtos.DeleteSpaceUserPermissionRequest) (*dtos.DeleteSpaceUserPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.space.delete_user_permission").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if err := ctrl.SpaceApp.DeleteSpaceUserPermission(spaceDto.DeleteSpaceUserPermissionInput{
		RequesterId:  authCtx.UserID,
		SpaceId:      req.SpaceId,
		TargetUserId: req.UserId,
	}); err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Space not found", Type: "SPACE_NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to delete permission")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete permission", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.DeleteSpaceUserPermissionResponse{Message: "permission deleted"}, nil
}
