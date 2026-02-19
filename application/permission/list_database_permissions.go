package permission

import (
	"fmt"

	databaseDto "github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// ListDatabasePermissions returns all permissions for a database.
// The requester must have access to the parent space.
func (app *PermissionApplication) ListDatabasePermissions(input databaseDto.ListDatabasePermissionsInput) (*databaseDto.ListDatabasePermissionsOutput, error) {
	dbResult, err := app.DatabaseApp.GetDatabaseById(databaseDto.GetDatabaseByIdInput{DatabaseId: input.DatabaseId})
	if err != nil {
		return nil, fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: dbResult.Database.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	perms, err := app.PermissionPers.ListByResource(domain.PermissionTypeDatabase, input.DatabaseId)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	output := &databaseDto.ListDatabasePermissionsOutput{
		Permissions: make([]databaseDto.DatabasePermissionItem, len(perms)),
	}

	for i, perm := range perms {
		item := databaseDto.DatabasePermissionItem{
			Id:   perm.Id,
			Role: string(perm.Role),
		}
		if perm.UserId != nil {
			item.UserId = perm.UserId
			if perm.User != nil {
				item.Username = &perm.User.Username
			}
		}
		if perm.GroupId != nil {
			item.GroupId = perm.GroupId
			if perm.Group != nil {
				item.GroupName = &perm.Group.Name
			}
		}
		output.Permissions[i] = item
	}

	return output, nil
}
