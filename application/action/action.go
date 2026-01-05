package action

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/action/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type ActionApp struct {
	Config        config.Config
	Logger        zerolog.Logger
	ActionPers    domain.ActionPers
	ActionRunPers domain.ActionRunPers
}

func NewActionApp(config config.Config, logger zerolog.Logger, actionPers domain.ActionPers, actionRunPers domain.ActionRunPers) *ActionApp {
	return &ActionApp{
		Config:        config,
		Logger:        logger,
		ActionPers:    actionPers,
		ActionRunPers: actionRunPers,
	}
}

func (app *ActionApp) CreateAction(input dto.CreateActionInput) (*dto.CreateActionOutput, error) {
	// Build steps JSONB
	stepsJSON, err := json.Marshal(input.Steps)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal steps: %w", err)
	}
	var steps domain.JSONB
	json.Unmarshal(stepsJSON, &steps)

	// Build trigger config JSONB
	var triggerConfig domain.JSONB
	if input.TriggerConfig != nil {
		triggerConfig = domain.JSONB(input.TriggerConfig)
	}

	action := &domain.Action{
		Id:            uuid.New().String(),
		UserId:        input.UserId,
		SpaceId:       input.SpaceId,
		DatabaseId:    input.DatabaseId,
		Name:          input.Name,
		Description:   input.Description,
		TriggerType:   domain.ActionTriggerType(input.TriggerType),
		TriggerConfig: triggerConfig,
		Steps:         steps,
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := app.ActionPers.Create(action); err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	return &dto.CreateActionOutput{
		Id:          action.Id,
		Name:        action.Name,
		Description: action.Description,
		TriggerType: input.TriggerType,
		Active:      action.Active,
		CreatedAt:   action.CreatedAt,
	}, nil
}

func (app *ActionApp) ListActions(input dto.ListActionsInput) (*dto.ListActionsOutput, error) {
	actions, err := app.ActionPers.GetByUserId(input.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to list actions: %w", err)
	}

	output := &dto.ListActionsOutput{
		Actions: make([]dto.ActionItem, len(actions)),
	}

	for i, a := range actions {
		item := dto.ActionItem{
			Id:           a.Id,
			Name:         a.Name,
			Description:  a.Description,
			SpaceId:      a.SpaceId,
			DatabaseId:   a.DatabaseId,
			TriggerType:  string(a.TriggerType),
			Active:       a.Active,
			LastRunAt:    a.LastRunAt,
			LastError:    a.LastError,
			RunCount:     a.RunCount,
			SuccessCount: a.SuccessCount,
			FailureCount: a.FailureCount,
			CreatedAt:    a.CreatedAt,
		}

		if a.Space != nil {
			item.SpaceName = &a.Space.Name
		}

		output.Actions[i] = item
	}

	return output, nil
}

func (app *ActionApp) GetAction(input dto.GetActionInput) (*dto.GetActionOutput, error) {
	action, err := app.ActionPers.GetById(input.ActionId)
	if err != nil {
		return nil, fmt.Errorf("action not found: %w", err)
	}

	if action.UserId != input.UserId {
		return nil, fmt.Errorf("access denied")
	}

	// Parse steps
	var steps []dto.ActionStep
	if action.Steps != nil {
		stepsJSON, _ := json.Marshal(action.Steps)
		json.Unmarshal(stepsJSON, &steps)
	}

	output := &dto.GetActionOutput{
		Id:            action.Id,
		Name:          action.Name,
		Description:   action.Description,
		SpaceId:       action.SpaceId,
		DatabaseId:    action.DatabaseId,
		TriggerType:   string(action.TriggerType),
		TriggerConfig: map[string]interface{}(action.TriggerConfig),
		Steps:         steps,
		Active:        action.Active,
		LastRunAt:     action.LastRunAt,
		LastError:     action.LastError,
		RunCount:      action.RunCount,
		SuccessCount:  action.SuccessCount,
		FailureCount:  action.FailureCount,
		CreatedAt:     action.CreatedAt,
		UpdatedAt:     action.UpdatedAt,
	}

	if action.Space != nil {
		output.SpaceName = &action.Space.Name
	}

	return output, nil
}

func (app *ActionApp) UpdateAction(input dto.UpdateActionInput) error {
	action, err := app.ActionPers.GetById(input.ActionId)
	if err != nil {
		return fmt.Errorf("action not found: %w", err)
	}

	if action.UserId != input.UserId {
		return fmt.Errorf("access denied")
	}

	if input.Name != nil {
		action.Name = *input.Name
	}

	if input.Description != nil {
		action.Description = *input.Description
	}

	if input.TriggerType != nil {
		action.TriggerType = domain.ActionTriggerType(*input.TriggerType)
	}

	if input.TriggerConfig != nil {
		action.TriggerConfig = domain.JSONB(input.TriggerConfig)
	}

	if input.Steps != nil {
		stepsJSON, err := json.Marshal(input.Steps)
		if err != nil {
			return fmt.Errorf("failed to marshal steps: %w", err)
		}
		var steps domain.JSONB
		json.Unmarshal(stepsJSON, &steps)
		action.Steps = steps
	}

	if input.Active != nil {
		action.Active = *input.Active
	}

	action.UpdatedAt = time.Now()

	if err := app.ActionPers.Update(action); err != nil {
		return fmt.Errorf("failed to update action: %w", err)
	}

	return nil
}

func (app *ActionApp) DeleteAction(input dto.DeleteActionInput) error {
	action, err := app.ActionPers.GetById(input.ActionId)
	if err != nil {
		return fmt.Errorf("action not found: %w", err)
	}

	if action.UserId != input.UserId {
		return fmt.Errorf("access denied")
	}

	if err := app.ActionPers.Delete(input.ActionId); err != nil {
		return fmt.Errorf("failed to delete action: %w", err)
	}

	return nil
}

func (app *ActionApp) GetRuns(input dto.GetRunsInput) (*dto.GetRunsOutput, error) {
	// Verify ownership
	action, err := app.ActionPers.GetById(input.ActionId)
	if err != nil {
		return nil, fmt.Errorf("action not found: %w", err)
	}

	if action.UserId != input.UserId {
		return nil, fmt.Errorf("access denied")
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	runs, err := app.ActionRunPers.GetByActionId(input.ActionId, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get runs: %w", err)
	}

	output := &dto.GetRunsOutput{
		Runs: make([]dto.RunItem, len(runs)),
	}

	for i, r := range runs {
		output.Runs[i] = dto.RunItem{
			Id:        r.Id,
			Success:   r.Success,
			Error:     r.Error,
			Duration:  r.Duration,
			CreatedAt: r.CreatedAt,
		}
	}

	return output, nil
}

// ExecuteActions triggers all matching actions for an event
func (app *ActionApp) ExecuteActions(input dto.ExecuteActionInput) {
	logger := app.Logger.With().Str("component", "action.execute").Str("trigger", input.TriggerType).Logger()

	actions, err := app.ActionPers.GetActiveByTrigger(domain.ActionTriggerType(input.TriggerType), input.SpaceId, input.DatabaseId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get actions for trigger")
		return
	}

	for _, action := range actions {
		go app.executeAction(action, input.TriggerData)
	}
}

func (app *ActionApp) executeAction(action domain.Action, triggerData map[string]interface{}) {
	logger := app.Logger.With().
		Str("component", "action.execute").
		Str("action_id", action.Id).
		Str("trigger", string(action.TriggerType)).
		Logger()

	start := time.Now()

	// Update last run
	_ = app.ActionPers.UpdateLastRun(action.Id)

	// Parse steps
	var steps []dto.ActionStep
	if action.Steps != nil {
		stepsJSON, _ := json.Marshal(action.Steps)
		json.Unmarshal(stepsJSON, &steps)
	}

	// Execute each step
	stepsResult := make([]map[string]interface{}, len(steps))
	var execError error

	for i, step := range steps {
		result, err := app.executeStep(step, triggerData)
		stepsResult[i] = map[string]interface{}{
			"step":    i + 1,
			"type":    step.Type,
			"success": err == nil,
			"result":  result,
		}
		if err != nil {
			stepsResult[i]["error"] = err.Error()
			execError = err
			break // Stop on first error
		}
	}

	duration := int(time.Since(start).Milliseconds())

	// Record the run
	run := &domain.ActionRun{
		Id:          uuid.New().String(),
		ActionId:    action.Id,
		TriggerData: domain.JSONB(triggerData),
		StepsResult: domain.JSONB{"steps": stepsResult},
		Success:     execError == nil,
		Duration:    duration,
		CreatedAt:   time.Now(),
	}

	if execError != nil {
		run.Error = execError.Error()
		_ = app.ActionPers.RecordFailure(action.Id, execError.Error())
		logger.Warn().Err(execError).Msg("action execution failed")
	} else {
		_ = app.ActionPers.IncrementSuccess(action.Id)
		logger.Debug().Msg("action executed successfully")
	}

	if err := app.ActionRunPers.Create(run); err != nil {
		logger.Error().Err(err).Msg("failed to record action run")
	}
}

func (app *ActionApp) executeStep(step dto.ActionStep, triggerData map[string]interface{}) (interface{}, error) {
	// This is a simplified implementation - in production, you'd have actual step executors
	switch domain.ActionStepType(step.Type) {
	case domain.StepSendWebhook:
		// Would call webhook service here
		return map[string]string{"status": "webhook_sent"}, nil
	case domain.StepSendEmail:
		// Would call email service here
		return map[string]string{"status": "email_sent"}, nil
	case domain.StepSendSlack:
		// Would call Slack API here
		return map[string]string{"status": "slack_sent"}, nil
	case domain.StepUpdateProperty:
		// Would update property in database
		return map[string]string{"status": "property_updated"}, nil
	case domain.StepAddComment:
		// Would add comment
		return map[string]string{"status": "comment_added"}, nil
	default:
		return nil, fmt.Errorf("unsupported step type: %s", step.Type)
	}
}
