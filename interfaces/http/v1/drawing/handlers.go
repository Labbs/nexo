package drawing

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	drawingDto "github.com/labbs/nexo/application/drawing/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/interfaces/http/v1/drawing/dtos"
)

func (ctrl *Controller) CreateDrawing(ctx *fiber.Ctx, req dtos.CreateDrawingRequest) (*dtos.CreateDrawingResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.drawing.create").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if req.Name == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Name is required", Type: "BAD_REQUEST"}
	}

	if req.SpaceId == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Space ID is required", Type: "BAD_REQUEST"}
	}

	result, err := ctrl.DrawingApp.CreateDrawing(drawingDto.CreateDrawingInput{
		UserId:     authCtx.UserID,
		SpaceId:    req.SpaceId,
		DocumentId: req.DocumentId,
		Name:       req.Name,
		Icon:       req.Icon,
		Elements:   req.Elements,
		AppState:   req.AppState,
		Files:      req.Files,
		Thumbnail:  req.Thumbnail,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		logger.Error().Err(err).Msg("failed to create drawing")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create drawing", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.CreateDrawingResponse{
		Id:        result.Id,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
	}, nil
}

func (ctrl *Controller) ListDrawings(ctx *fiber.Ctx, req dtos.ListDrawingsRequest) (*dtos.ListDrawingsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.drawing.list").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if req.SpaceId == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Space ID is required", Type: "BAD_REQUEST"}
	}

	result, err := ctrl.DrawingApp.ListDrawings(drawingDto.ListDrawingsInput{
		UserId:  authCtx.UserID,
		SpaceId: req.SpaceId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		logger.Error().Err(err).Msg("failed to list drawings")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list drawings", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.ListDrawingsResponse{Drawings: make([]dtos.DrawingItem, len(result.Drawings))}
	for i, d := range result.Drawings {
		resp.Drawings[i] = dtos.DrawingItem{
			Id:         d.Id,
			DocumentId: d.DocumentId,
			Name:       d.Name,
			Icon:       d.Icon,
			Thumbnail:  d.Thumbnail,
			CreatedBy:  d.CreatedBy,
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) GetDrawing(ctx *fiber.Ctx, req dtos.GetDrawingRequest) (*dtos.GetDrawingResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.drawing.get").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DrawingApp.GetDrawing(drawingDto.GetDrawingInput{
		UserId:    authCtx.UserID,
		DrawingId: req.DrawingId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Drawing not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to get drawing")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get drawing", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.GetDrawingResponse{
		Id:         result.Id,
		SpaceId:    result.SpaceId,
		DocumentId: result.DocumentId,
		Name:       result.Name,
		Icon:       result.Icon,
		Elements:   result.Elements,
		AppState:   result.AppState,
		Files:      result.Files,
		Thumbnail:  result.Thumbnail,
		CreatedBy:  result.CreatedBy,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}, nil
}

func (ctrl *Controller) UpdateDrawing(ctx *fiber.Ctx, req dtos.UpdateDrawingRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.drawing.update").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DrawingApp.UpdateDrawing(drawingDto.UpdateDrawingInput{
		UserId:    authCtx.UserID,
		DrawingId: req.DrawingId,
		Name:      req.Name,
		Icon:      req.Icon,
		Elements:  req.Elements,
		AppState:  req.AppState,
		Files:     req.Files,
		Thumbnail: req.Thumbnail,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Drawing not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update drawing")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update drawing", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Drawing updated successfully"}, nil
}

func (ctrl *Controller) DeleteDrawing(ctx *fiber.Ctx, req dtos.DeleteDrawingRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.drawing.delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DrawingApp.DeleteDrawing(drawingDto.DeleteDrawingInput{
		UserId:    authCtx.UserID,
		DrawingId: req.DrawingId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Drawing not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to delete drawing")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete drawing", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Drawing deleted successfully"}, nil
}

// Permission handlers

func (ctrl *Controller) ListDrawingPermissions(ctx *fiber.Ctx, req dtos.ListDrawingPermissionsRequest) (*dtos.ListDrawingPermissionsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.drawing.list_permissions").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DrawingApp.ListDrawingPermissions(drawingDto.ListDrawingPermissionsInput{
		RequesterId: authCtx.UserID,
		DrawingId:   req.DrawingId,
	})
	if err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Drawing not found", Type: "NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to list drawing permissions")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list permissions", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	resp := &dtos.ListDrawingPermissionsResponse{Permissions: make([]dtos.DrawingPermission, len(result.Permissions))}
	for i, p := range result.Permissions {
		perm := dtos.DrawingPermission{
			Id:     p.Id,
			UserId: p.UserId,
			Role:   string(p.Role),
		}
		if p.User != nil {
			perm.Username = &p.User.Username
		}
		resp.Permissions[i] = perm
	}
	return resp, nil
}

func (ctrl *Controller) UpsertDrawingUserPermission(ctx *fiber.Ctx, req dtos.UpsertDrawingUserPermissionRequest) (*dtos.UpsertDrawingUserPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.drawing.upsert_permission").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	var role domain.PermissionRole
	switch req.Role {
	case string(domain.PermissionRoleOwner):
		role = domain.PermissionRoleOwner
	case string(domain.PermissionRoleEditor):
		role = domain.PermissionRoleEditor
	case string(domain.PermissionRoleDenied):
		role = domain.PermissionRoleDenied
	default:
		role = domain.PermissionRoleViewer
	}

	if err := ctrl.DrawingApp.UpsertDrawingUserPermission(drawingDto.UpsertDrawingUserPermissionInput{
		RequesterId:  authCtx.UserID,
		DrawingId:    req.DrawingId,
		TargetUserId: req.UserId,
		Role:         role,
	}); err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Drawing not found", Type: "NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to upsert drawing permission")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to upsert permission", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.UpsertDrawingUserPermissionResponse{Message: "permission updated"}, nil
}

func (ctrl *Controller) DeleteDrawingUserPermission(ctx *fiber.Ctx, req dtos.DeleteDrawingUserPermissionRequest) (*dtos.DeleteDrawingUserPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.drawing.delete_permission").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if err := ctrl.DrawingApp.DeleteDrawingUserPermission(drawingDto.DeleteDrawingUserPermissionInput{
		RequesterId:  authCtx.UserID,
		DrawingId:    req.DrawingId,
		TargetUserId: req.UserId,
	}); err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Drawing not found", Type: "NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to delete drawing permission")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete permission", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.DeleteDrawingUserPermissionResponse{Message: "permission deleted"}, nil
}
