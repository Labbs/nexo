package action

import (
	"fmt"

	"github.com/labbs/nexo/application/action/dto"
)

func (app *ActionApplication) DeleteAction(input dto.DeleteActionInput) error {
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
