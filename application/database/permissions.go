package database

import (
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
	"github.com/labbs/nexo/domain"
)

// ListDatabasePermissions returns all permissions for a database
func (app *DatabaseApp) ListDatabasePermissions(input dto.ListDatabasePermissionsInput) (*dto.ListDatabasePermissionsOutput, error) {
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

// UpsertDatabasePermission adds or updates a permission for a database
func (app *DatabaseApp) UpsertDatabasePermission(input dto.UpsertDatabasePermissionInput) error {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Verify user has permission to manage permissions (creator or space admin)
	space, err := app.SpacePers.GetSpaceById(database.SpaceId)
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	spaceRole := space.GetUserRole(input.UserId)
	if spaceRole == nil {
		return fmt.Errorf("access denied")
	}

	// Only creator or space admin/owner can manage permissions
	isCreator := database.CreatedBy == input.UserId
	isSpaceAdmin := *spaceRole == domain.PermissionRoleOwner || *spaceRole == domain.PermissionRoleAdmin
	if !isCreator && !isSpaceAdmin {
		return fmt.Errorf("only creator or space admins can manage permissions")
	}

	// Validate role
	role := domain.PermissionRole(input.Role)
	if role != domain.PermissionRoleEditor && role != domain.PermissionRoleViewer && role != domain.PermissionRoleDenied {
		return fmt.Errorf("invalid role: %s", input.Role)
	}

	// Validate that either UserId or GroupId is provided, but not both
	if (input.TargetUserId == nil && input.GroupId == nil) || (input.TargetUserId != nil && input.GroupId != nil) {
		return fmt.Errorf("must provide either user_id or group_id, but not both")
	}

	if input.TargetUserId != nil {
		if err := app.PermissionPers.UpsertUser(domain.PermissionTypeDatabase, input.DatabaseId, *input.TargetUserId, role); err != nil {
			return fmt.Errorf("failed to upsert permission: %w", err)
		}
	} else {
		if err := app.PermissionPers.UpsertGroup(domain.PermissionTypeDatabase, input.DatabaseId, *input.GroupId, role); err != nil {
			return fmt.Errorf("failed to upsert permission: %w", err)
		}
	}

	return nil
}

// DeleteDatabasePermission removes a permission from a database
func (app *DatabaseApp) DeleteDatabasePermission(input dto.DeleteDatabasePermissionInput) error {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Verify user has permission to manage permissions (creator or space admin)
	space, err := app.SpacePers.GetSpaceById(database.SpaceId)
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	spaceRole := space.GetUserRole(input.UserId)
	if spaceRole == nil {
		return fmt.Errorf("access denied")
	}

	// Only creator or space admin/owner can manage permissions
	isCreator := database.CreatedBy == input.UserId
	isSpaceAdmin := *spaceRole == domain.PermissionRoleOwner || *spaceRole == domain.PermissionRoleAdmin
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
