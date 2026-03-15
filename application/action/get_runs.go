package action

import (
	"fmt"

	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/action/dto"
)

func (app *ActionApplication) GetRuns(input dto.GetRunsInput) (*dto.GetRunsOutput, error) {
	// Verify ownership
	action, err := app.ActionPers.GetById(input.ActionId)
	if err != nil {
		return nil, fmt.Errorf("action not found: %w", err)
	}

	if action.UserId != input.UserId {
		return nil, apperrors.ErrAccessDenied
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
