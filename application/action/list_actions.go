package action

import (
	"fmt"

	"github.com/labbs/nexo/application/action/dto"
)

func (app *ActionApplication) ListActions(input dto.ListActionsInput) (*dto.ListActionsOutput, error) {
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
