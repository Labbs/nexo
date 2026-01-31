package apikey

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	apikeyDto "github.com/labbs/nexo/application/apikey/dto"
	"github.com/labbs/nexo/interfaces/http/v1/apikey/dtos"
)

func (ctrl *Controller) ListApiKeys(ctx *fiber.Ctx, _ dtos.EmptyRequest) (*dtos.ListApiKeysResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.apikey.list").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.ApiKeyApp.ListApiKeys(apikeyDto.ListApiKeysInput{
		UserId: authCtx.UserID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to list API keys")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list API keys", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.ListApiKeysResponse{ApiKeys: make([]dtos.ApiKeyItem, len(result.ApiKeys))}
	for i, k := range result.ApiKeys {
		resp.ApiKeys[i] = dtos.ApiKeyItem{
			Id:         k.Id,
			Name:       k.Name,
			KeyPrefix:  k.KeyPrefix,
			Scopes:     k.Scopes,
			LastUsedAt: k.LastUsedAt,
			ExpiresAt:  k.ExpiresAt,
			CreatedAt:  k.CreatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) CreateApiKey(ctx *fiber.Ctx, req dtos.CreateApiKeyRequest) (*dtos.CreateApiKeyResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.apikey.create").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if req.Name == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Name is required", Type: "BAD_REQUEST"}
	}

	if len(req.Scopes) == 0 {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "At least one scope is required", Type: "BAD_REQUEST"}
	}

	result, err := ctrl.ApiKeyApp.CreateApiKey(apikeyDto.CreateApiKeyInput{
		UserId:    authCtx.UserID,
		Name:      req.Name,
		Scopes:    req.Scopes,
		ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create API key")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create API key", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.CreateApiKeyResponse{
		Id:        result.Id,
		Name:      result.Name,
		Key:       result.Key,
		KeyPrefix: result.KeyPrefix,
		Scopes:    result.Scopes,
		ExpiresAt: result.ExpiresAt,
		CreatedAt: result.CreatedAt,
	}, nil
}

func (ctrl *Controller) UpdateApiKey(ctx *fiber.Ctx, req dtos.UpdateApiKeyRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.apikey.update").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.ApiKeyApp.UpdateApiKey(apikeyDto.UpdateApiKeyInput{
		UserId:   authCtx.UserID,
		ApiKeyId: req.ApiKeyId,
		Name:     req.Name,
		Scopes:   req.Scopes,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "API key not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update API key")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update API key", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "API key updated"}, nil
}

func (ctrl *Controller) DeleteApiKey(ctx *fiber.Ctx, req dtos.DeleteApiKeyRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.apikey.delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.ApiKeyApp.DeleteApiKey(apikeyDto.DeleteApiKeyInput{
		UserId:   authCtx.UserID,
		ApiKeyId: req.ApiKeyId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "API key not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to delete API key")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete API key", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "API key deleted"}, nil
}

func (ctrl *Controller) GetAvailableScopes(ctx *fiber.Ctx, _ dtos.EmptyRequest) (*dtos.AvailableScopesResponse, *fiberoapi.ErrorResponse) {
	scopes := []dtos.ScopeInfo{
		{Scope: "read:documents", Description: "Read access to documents"},
		{Scope: "write:documents", Description: "Write access to documents (create, update, delete)"},
		{Scope: "read:spaces", Description: "Read access to spaces"},
		{Scope: "write:spaces", Description: "Write access to spaces"},
		{Scope: "read:comments", Description: "Read access to comments"},
		{Scope: "write:comments", Description: "Write access to comments"},
		{Scope: "manage:webhooks", Description: "Manage webhooks"},
		{Scope: "manage:databases", Description: "Manage databases"},
	}

	return &dtos.AvailableScopesResponse{Scopes: scopes}, nil
}
