package permission

import (
	"fmt"

	"github.com/labbs/nexo/application/permission/dto"
	"github.com/labbs/nexo/domain"
)

// AssignOwnerPermission assigns a permission to a resource creator without authorization checks.
// This is used internally during resource creation to auto-grant the creator access.
func (app *PermissionApplication) AssignOwnerPermission(input dto.AssignOwnerPermissionInput) error {
	if err := app.PermissionPers.UpsertUser(domain.PermissionType(input.ResourceType), input.ResourceId, input.UserId, domain.PermissionRole(input.Role)); err != nil {
		return fmt.Errorf("failed to assign owner permission: %w", err)
	}

	return nil
}
