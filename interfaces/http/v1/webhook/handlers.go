package webhook

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	webhookDto "github.com/labbs/nexo/application/webhook/dto"
	"github.com/labbs/nexo/interfaces/http/v1/webhook/dtos"
)

func (ctrl *Controller) ListWebhooks(ctx *fiber.Ctx, _ dtos.EmptyRequest) (*dtos.ListWebhooksResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.webhook.list").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.WebhookApp.ListWebhooks(webhookDto.ListWebhooksInput{
		UserId: authCtx.UserID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to list webhooks")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list webhooks", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.ListWebhooksResponse{Webhooks: make([]dtos.WebhookItem, len(result.Webhooks))}
	for i, w := range result.Webhooks {
		resp.Webhooks[i] = dtos.WebhookItem{
			Id:           w.Id,
			Name:         w.Name,
			Url:          w.Url,
			SpaceId:      w.SpaceId,
			SpaceName:    w.SpaceName,
			Events:       w.Events,
			Active:       w.Active,
			LastError:    w.LastError,
			LastErrorAt:  w.LastErrorAt,
			SuccessCount: w.SuccessCount,
			FailureCount: w.FailureCount,
			CreatedAt:    w.CreatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) CreateWebhook(ctx *fiber.Ctx, req dtos.CreateWebhookRequest) (*dtos.CreateWebhookResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.webhook.create").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if req.Name == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Name is required", Type: "BAD_REQUEST"}
	}

	if req.Url == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "URL is required", Type: "BAD_REQUEST"}
	}

	if len(req.Events) == 0 {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "At least one event is required", Type: "BAD_REQUEST"}
	}

	result, err := ctrl.WebhookApp.CreateWebhook(webhookDto.CreateWebhookInput{
		UserId:  authCtx.UserID,
		SpaceId: req.SpaceId,
		Name:    req.Name,
		Url:     req.Url,
		Events:  req.Events,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create webhook")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create webhook", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.CreateWebhookResponse{
		Id:     result.Id,
		Name:   result.Name,
		Url:    result.Url,
		Secret: result.Secret,
		Events: result.Events,
		Active: result.Active,
	}, nil
}

func (ctrl *Controller) GetWebhook(ctx *fiber.Ctx, req dtos.GetWebhookRequest) (*dtos.GetWebhookResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.webhook.get").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.WebhookApp.GetWebhook(webhookDto.GetWebhookInput{
		UserId:    authCtx.UserID,
		WebhookId: req.WebhookId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Webhook not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to get webhook")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get webhook", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.GetWebhookResponse{
		Id:           result.Id,
		Name:         result.Name,
		Url:          result.Url,
		Secret:       result.Secret,
		SpaceId:      result.SpaceId,
		SpaceName:    result.SpaceName,
		Events:       result.Events,
		Active:       result.Active,
		LastError:    result.LastError,
		LastErrorAt:  result.LastErrorAt,
		SuccessCount: result.SuccessCount,
		FailureCount: result.FailureCount,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}, nil
}

func (ctrl *Controller) UpdateWebhook(ctx *fiber.Ctx, req dtos.UpdateWebhookRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.webhook.update").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.WebhookApp.UpdateWebhook(webhookDto.UpdateWebhookInput{
		UserId:    authCtx.UserID,
		WebhookId: req.WebhookId,
		Name:      req.Name,
		Url:       req.Url,
		Events:    req.Events,
		Active:    req.Active,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Webhook not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update webhook")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update webhook", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Webhook updated"}, nil
}

func (ctrl *Controller) DeleteWebhook(ctx *fiber.Ctx, req dtos.DeleteWebhookRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.webhook.delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.WebhookApp.DeleteWebhook(webhookDto.DeleteWebhookInput{
		UserId:    authCtx.UserID,
		WebhookId: req.WebhookId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Webhook not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to delete webhook")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete webhook", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Webhook deleted"}, nil
}

func (ctrl *Controller) GetDeliveries(ctx *fiber.Ctx, req dtos.GetDeliveriesRequest) (*dtos.GetDeliveriesResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.webhook.deliveries").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	result, err := ctrl.WebhookApp.GetDeliveries(webhookDto.GetDeliveriesInput{
		UserId:    authCtx.UserID,
		WebhookId: req.WebhookId,
		Limit:     limit,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Webhook not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to get deliveries")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get deliveries", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.GetDeliveriesResponse{Deliveries: make([]dtos.DeliveryItem, len(result.Deliveries))}
	for i, d := range result.Deliveries {
		resp.Deliveries[i] = dtos.DeliveryItem{
			Id:         d.Id,
			Event:      d.Event,
			StatusCode: d.StatusCode,
			Success:    d.Success,
			Duration:   d.Duration,
			CreatedAt:  d.CreatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) GetAvailableEvents(ctx *fiber.Ctx, _ dtos.EmptyRequest) (*dtos.AvailableEventsResponse, *fiberoapi.ErrorResponse) {
	events := []dtos.EventInfo{
		{Event: "document.created", Description: "Triggered when a document is created"},
		{Event: "document.updated", Description: "Triggered when a document is updated"},
		{Event: "document.deleted", Description: "Triggered when a document is deleted"},
		{Event: "comment.created", Description: "Triggered when a comment is created"},
		{Event: "comment.resolved", Description: "Triggered when a comment is resolved"},
		{Event: "space.created", Description: "Triggered when a space is created"},
		{Event: "space.updated", Description: "Triggered when a space is updated"},
	}

	return &dtos.AvailableEventsResponse{Events: events}, nil
}
