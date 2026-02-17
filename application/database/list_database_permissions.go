package database

import (
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
	"github.com/labbs/nexo/domain"
)

// ListDatabasePermissions returns all permissions for a database
func (app *DatabaseApplication) ListDatabasePermissions(input dto.ListDatabasePermissionsInput) (*dto.ListDatabasePermissionsOutput, error) {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return nil, fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	space, err := app.SpacePers.GetSpaceById(database.SpaceId)
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	perms, err := app.PermissionPers.ListByResource(domain.PermissionTypeDatabase, input.DatabaseId)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	output := &dto.ListDatabasePermissionsOutput{
		Permissions: make([]dto.DatabasePermissionItem, len(perms)),
	}

	for i, perm := range perms {
		item := dto.DatabasePermissionItem{
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
