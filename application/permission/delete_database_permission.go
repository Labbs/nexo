package permission

import (
	"fmt"

	databaseDto "github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// DeleteDatabasePermission removes a permission from a database.
// Only the database creator or a space admin/owner can manage permissions.
// Supports both user and group permissions (mutually exclusive).
func (app *PermissionApplication) DeleteDatabasePermission(input databaseDto.DeleteDatabasePermissionInput) error {
	dbResult, err := app.DatabaseApp.GetDatabaseById(databaseDto.GetDatabaseByIdInput{DatabaseId: input.DatabaseId})
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Verify user has permission to manage permissions (creator or space admin)
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: dbResult.Database.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	spaceRole := spaceResult.Space.GetUserRole(input.UserId)
	if spaceRole == nil {
		return fmt.Errorf("access denied")
	}

	// Only creator or space admin/owner can manage permissions
	isCreator := dbResult.Database.CreatedBy == input.UserId
	isSpaceAdmin := *spaceRole == "owner" || *spaceRole == "admin"
	if !isCreator && !isSpaceAdmin {
		return fmt.Errorf("only creator or space admins can manage permissions")
	}

	// Validate that either UserId or GroupId is provided, but not both
	if (input.TargetUserId == nil && input.GroupId == nil) || (input.TargetUserId != nil && input.GroupId != nil) {
		return fmt.Errorf("must provide either user_id or group_id, but not both")
	}

	if input.TargetUserId != nil {
		if err := app.PermissionPers.DeleteUser(domain.PermissionTypeDatabase, input.DatabaseId, *input.TargetUserId); err != nil {
			return fmt.Errorf("failed to delete permission: %w", err)
		}
	} else {
		if err := app.PermissionPers.DeleteGroup(domain.PermissionTypeDatabase, input.DatabaseId, *input.GroupId); err != nil {
			return fmt.Errorf("failed to delete permission: %w", err)
		}
	}

	return nil
}
