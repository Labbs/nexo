package action

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	actionDto "github.com/labbs/nexo/application/action/dto"
	"github.com/labbs/nexo/interfaces/http/v1/action/dtos"
)

func (ctrl *Controller) ListActions(ctx *fiber.Ctx, _ dtos.EmptyRequest) (*dtos.ListActionsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.action.list").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.ActionApp.ListActions(actionDto.ListActionsInput{
		UserId: authCtx.UserID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to list actions")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list actions", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.ListActionsResponse{Actions: make([]dtos.ActionItem, len(result.Actions))}
	for i, a := range result.Actions {
		resp.Actions[i] = dtos.ActionItem{
			Id:           a.Id,
			Name:         a.Name,
			Description:  a.Description,
			SpaceId:      a.SpaceId,
			SpaceName:    a.SpaceName,
			DatabaseId:   a.DatabaseId,
			TriggerType:  a.TriggerType,
			Active:       a.Active,
			LastRunAt:    a.LastRunAt,
			LastError:    a.LastError,
			RunCount:     a.RunCount,
			SuccessCount: a.SuccessCount,
			FailureCount: a.FailureCount,
			CreatedAt:    a.CreatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) CreateAction(ctx *fiber.Ctx, req dtos.CreateActionRequest) (*dtos.CreateActionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.action.create").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if req.Name == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Name is required", Type: "BAD_REQUEST"}
	}

	if req.TriggerType == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Trigger type is required", Type: "BAD_REQUEST"}
	}

	if len(req.Steps) == 0 {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "At least one step is required", Type: "BAD_REQUEST"}
	}

	// Convert steps
	steps := make([]actionDto.ActionStep, len(req.Steps))
	for i, s := range req.Steps {
		steps[i] = actionDto.ActionStep{
			Type:   s.Type,
			Config: s.Config,
		}
	}

	result, err := ctrl.ActionApp.CreateAction(actionDto.CreateActionInput{
		UserId:        authCtx.UserID,
		SpaceId:       req.SpaceId,
		DatabaseId:    req.DatabaseId,
		Name:          req.Name,
		Description:   req.Description,
		TriggerType:   req.TriggerType,
		TriggerConfig: req.TriggerConfig,
		Steps:         steps,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create action")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create action", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.CreateActionResponse{
		Id:          result.Id,
		Name:        result.Name,
		Description: result.Description,
		TriggerType: result.TriggerType,
		Active:      result.Active,
		CreatedAt:   result.CreatedAt,
	}, nil
}

func (ctrl *Controller) GetAction(ctx *fiber.Ctx, req dtos.GetActionRequest) (*dtos.GetActionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.action.get").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.ActionApp.GetAction(actionDto.GetActionInput{
		UserId:   authCtx.UserID,
		ActionId: req.ActionId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Action not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to get action")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get action", Type: "INTERNAL_SERVER_ERROR"}
	}

	steps := make([]dtos.ActionStep, len(result.Steps))
	for i, s := range result.Steps {
		steps[i] = dtos.ActionStep{
			Type:   s.Type,
			Config: s.Config,
		}
	}

	return &dtos.GetActionResponse{
		Id:            result.Id,
		Name:          result.Name,
		Description:   result.Description,
		SpaceId:       result.SpaceId,
		SpaceName:     result.SpaceName,
		DatabaseId:    result.DatabaseId,
		TriggerType:   result.TriggerType,
		TriggerConfig: result.TriggerConfig,
		Steps:         steps,
		Active:        result.Active,
		LastRunAt:     result.LastRunAt,
		LastError:     result.LastError,
		RunCount:      result.RunCount,
		SuccessCount:  result.SuccessCount,
		FailureCount:  result.FailureCount,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}, nil
}

func (ctrl *Controller) UpdateAction(ctx *fiber.Ctx, req dtos.UpdateActionRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.action.update").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	var steps *[]actionDto.ActionStep
	if req.Steps != nil {
		s := make([]actionDto.ActionStep, len(*req.Steps))
		for i, step := range *req.Steps {
			s[i] = actionDto.ActionStep{
				Type:   step.Type,
				Config: step.Config,
			}
		}
		steps = &s
	}

	err = ctrl.ActionApp.UpdateAction(actionDto.UpdateActionInput{
		UserId:        authCtx.UserID,
		ActionId:      req.ActionId,
		Name:          req.Name,
		Description:   req.Description,
		TriggerType:   req.TriggerType,
		TriggerConfig: req.TriggerConfig,
		Steps:         steps,
		Active:        req.Active,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Action not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update action")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update action", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Action updated"}, nil
}

func (ctrl *Controller) DeleteAction(ctx *fiber.Ctx, req dtos.DeleteActionRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.action.delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.ActionApp.DeleteAction(actionDto.DeleteActionInput{
		UserId:   authCtx.UserID,
		ActionId: req.ActionId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Action not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to delete action")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete action", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Action deleted"}, nil
}

func (ctrl *Controller) GetRuns(ctx *fiber.Ctx, req dtos.GetRunsRequest) (*dtos.GetRunsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.action.runs").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	result, err := ctrl.ActionApp.GetRuns(actionDto.GetRunsInput{
		UserId:   authCtx.UserID,
		ActionId: req.ActionId,
		Limit:    limit,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Action not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to get runs")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get runs", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.GetRunsResponse{Runs: make([]dtos.RunItem, len(result.Runs))}
	for i, r := range result.Runs {
		resp.Runs[i] = dtos.RunItem{
			Id:        r.Id,
			Success:   r.Success,
			Error:     r.Error,
			Duration:  r.Duration,
			CreatedAt: r.CreatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) GetAvailableTriggers(ctx *fiber.Ctx, _ dtos.EmptyRequest) (*dtos.AvailableTriggersResponse, *fiberoapi.ErrorResponse) {
	triggers := []dtos.TriggerInfo{
		// Document triggers
		{Type: "document.created", Description: "When a document is created", Category: "Documents"},
		{Type: "document.updated", Description: "When a document is updated", Category: "Documents"},
		{Type: "document.deleted", Description: "When a document is deleted", Category: "Documents"},
		{Type: "document.moved", Description: "When a document is moved", Category: "Documents"},
		{Type: "document.shared", Description: "When a document is shared", Category: "Documents"},
		// Database triggers
		{Type: "row.created", Description: "When a database row is created", Category: "Databases"},
		{Type: "row.updated", Description: "When a database row is updated", Category: "Databases"},
		{Type: "row.deleted", Description: "When a database row is deleted", Category: "Databases"},
		{Type: "property.changed", Description: "When a specific property value changes", Category: "Databases"},
		// Comment triggers
		{Type: "comment.created", Description: "When a comment is added", Category: "Comments"},
		{Type: "comment.resolved", Description: "When a comment is resolved", Category: "Comments"},
		// Schedule triggers
		{Type: "schedule", Description: "Run on a schedule (cron)", Category: "Schedule"},
	}

	return &dtos.AvailableTriggersResponse{Triggers: triggers}, nil
}

func (ctrl *Controller) GetAvailableSteps(ctx *fiber.Ctx, _ dtos.EmptyRequest) (*dtos.AvailableStepsResponse, *fiberoapi.ErrorResponse) {
	steps := []dtos.StepInfo{
		// Notification steps
		{Type: "send_email", Description: "Send an email notification", Category: "Notifications"},
		{Type: "send_slack", Description: "Send a Slack message", Category: "Notifications"},
		{Type: "send_webhook", Description: "Send a webhook request", Category: "Notifications"},
		// Document steps
		{Type: "create_document", Description: "Create a new document", Category: "Documents"},
		{Type: "update_document", Description: "Update a document", Category: "Documents"},
		{Type: "move_document", Description: "Move a document", Category: "Documents"},
		{Type: "duplicate_document", Description: "Duplicate a document", Category: "Documents"},
		// Database steps
		{Type: "create_row", Description: "Create a database row", Category: "Databases"},
		{Type: "update_row", Description: "Update a database row", Category: "Databases"},
		{Type: "delete_row", Description: "Delete a database row", Category: "Databases"},
		{Type: "update_property", Description: "Update a property value", Category: "Databases"},
		// Misc steps
		{Type: "add_comment", Description: "Add a comment", Category: "Other"},
		{Type: "assign_user", Description: "Assign a user", Category: "Other"},
		{Type: "set_reminder", Description: "Set a reminder", Category: "Other"},
	}

	return &dtos.AvailableStepsResponse{Steps: steps}, nil
}
