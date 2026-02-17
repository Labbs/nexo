package action

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/action/dto"
	"github.com/labbs/nexo/domain"
)

// ExecuteActions triggers all matching actions for an event
func (app *ActionApplication) ExecuteActions(input dto.ExecuteActionInput) {
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

func (app *ActionApplication) executeAction(action domain.Action, triggerData map[string]interface{}) {
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

func (app *ActionApplication) executeStep(step dto.ActionStep, triggerData map[string]interface{}) (interface{}, error) {
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
