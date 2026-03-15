package action

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labbs/nexo/application/action/dto"
	"github.com/labbs/nexo/domain"
)

func (app *ActionApplication) UpdateAction(input dto.UpdateActionInput) error {
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
