package action

import (
	"encoding/json"
	"fmt"

	"github.com/labbs/nexo/application/action/dto"
)

func (app *ActionApplication) GetAction(input dto.GetActionInput) (*dto.GetActionOutput, error) {
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
