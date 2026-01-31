package database

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	databaseDto "github.com/labbs/nexo/application/database/dto"
	"github.com/labbs/nexo/interfaces/http/v1/database/dtos"
)

func (ctrl *Controller) ListDatabasePermissions(ctx *fiber.Ctx, req dtos.ListDatabasePermissionsRequest) (*dtos.ListDatabasePermissionsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.permissions.list").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DatabaseApp.ListDatabasePermissions(databaseDto.ListDatabasePermissionsInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to list permissions")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list permissions", Type: "INTERNAL_SERVER_ERROR"}
	}

	permissions := make([]dtos.DatabasePermissionItem, len(result.Permissions))
	for i, p := range result.Permissions {
		permissions[i] = dtos.DatabasePermissionItem{
			Id:        p.Id,
			UserId:    p.UserId,
			Username:  p.Username,
			GroupId:   p.GroupId,
			GroupName: p.GroupName,
			Role:      p.Role,
		}
	}

	return &dtos.ListDatabasePermissionsResponse{
		Permissions: permissions,
	}, nil
}

func (ctrl *Controller) UpsertDatabasePermission(ctx *fiber.Ctx, req dtos.UpsertDatabasePermissionRequest) (*dtos.UpsertDatabasePermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.permissions.upsert").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DatabaseApp.UpsertDatabasePermission(databaseDto.UpsertDatabasePermissionInput{
		UserId:       authCtx.UserID,
		DatabaseId:   req.DatabaseId,
		TargetUserId: req.UserId,
		GroupId:      req.GroupId,
		Role:         req.Role,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") || strings.Contains(err.Error(), "only creator") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: err.Error(), Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		if strings.Contains(err.Error(), "invalid role") || strings.Contains(err.Error(), "must provide") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: err.Error(), Type: "BAD_REQUEST"}
		}
		logger.Error().Err(err).Msg("failed to upsert permission")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to upsert permission", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.UpsertDatabasePermissionResponse{
		Success: true,
	}, nil
}

func (ctrl *Controller) DeleteDatabasePermission(ctx *fiber.Ctx, req dtos.DeleteDatabasePermissionRequest) (*dtos.DeleteDatabasePermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.permissions.delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DatabaseApp.DeleteDatabasePermission(databaseDto.DeleteDatabasePermissionInput{
		UserId:       authCtx.UserID,
		DatabaseId:   req.DatabaseId,
		TargetUserId: req.UserId,
		GroupId:      req.GroupId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") || strings.Contains(err.Error(), "only creator") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: err.Error(), Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		if strings.Contains(err.Error(), "must provide") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: err.Error(), Type: "BAD_REQUEST"}
		}
		logger.Error().Err(err).Msg("failed to delete permission")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete permission", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.DeleteDatabasePermissionResponse{
		Success: true,
	}, nil
}
