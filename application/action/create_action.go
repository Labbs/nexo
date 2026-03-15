package action

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/action/dto"
	"github.com/labbs/nexo/domain"
)

func (app *ActionApplication) CreateAction(input dto.CreateActionInput) (*dto.CreateActionOutput, error) {
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
