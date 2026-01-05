package admin

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/apikey"
	"github.com/labbs/nexo/application/group"
	"github.com/labbs/nexo/application/space"
	"github.com/labbs/nexo/application/user"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config    config.Config
	Logger    zerolog.Logger
	FiberOapi *fiberoapi.OApiGroup
	UserApp   *user.UserApp
	SpaceApp  *space.SpaceApp
	ApiKeyApp *apikey.ApiKeyApp
	GroupApp  *group.GroupApp
}

func SetupAdminRouter(controller Controller) {
	// Users management
	fiberoapi.Get(controller.FiberOapi, "/users", controller.ListUsers, fiberoapi.OpenAPIOptions{
		Summary:     "List all users",
		Description: "Retrieve a paginated list of all users (admin only)",
		OperationID: "admin.listUsers",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Put(controller.FiberOapi, "/users/:user_id/role", controller.UpdateUserRole, fiberoapi.OpenAPIOptions{
		Summary:     "Update user role",
		Description: "Change the role of a user (admin only)",
		OperationID: "admin.updateUserRole",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Put(controller.FiberOapi, "/users/:user_id/active", controller.UpdateUserActive, fiberoapi.OpenAPIOptions{
		Summary:     "Update user active status",
		Description: "Enable or disable a user account (admin only)",
		OperationID: "admin.updateUserActive",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/users/:user_id", controller.DeleteUser, fiberoapi.OpenAPIOptions{
		Summary:     "Delete user",
		Description: "Delete a user account (admin only)",
		OperationID: "admin.deleteUser",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Post(controller.FiberOapi, "/users/invite", controller.InviteUser, fiberoapi.OpenAPIOptions{
		Summary:     "Invite user",
		Description: "Invite a new user by email (admin only)",
		OperationID: "admin.inviteUser",
		Tags:        []string{"Admin"},
	})

	// Spaces management
	fiberoapi.Get(controller.FiberOapi, "/spaces", controller.ListAllSpaces, fiberoapi.OpenAPIOptions{
		Summary:     "List all spaces",
		Description: "Retrieve a paginated list of all spaces (admin only)",
		OperationID: "admin.listAllSpaces",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Post(controller.FiberOapi, "/spaces", controller.AdminCreateSpace, fiberoapi.OpenAPIOptions{
		Summary:     "Create space",
		Description: "Create a new space (admin only)",
		OperationID: "admin.createSpace",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Put(controller.FiberOapi, "/spaces/:space_id", controller.AdminUpdateSpace, fiberoapi.OpenAPIOptions{
		Summary:     "Update space",
		Description: "Update a space (admin only)",
		OperationID: "admin.updateSpace",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/spaces/:space_id", controller.AdminDeleteSpace, fiberoapi.OpenAPIOptions{
		Summary:     "Delete space",
		Description: "Delete a space (admin only)",
		OperationID: "admin.deleteSpace",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Get(controller.FiberOapi, "/spaces/:space_id/permissions", controller.AdminListSpacePermissions, fiberoapi.OpenAPIOptions{
		Summary:     "List space permissions",
		Description: "List all permissions for a space (admin only)",
		OperationID: "admin.listSpacePermissions",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Post(controller.FiberOapi, "/spaces/:space_id/permissions/users", controller.AdminAddSpaceUserPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Add user permission",
		Description: "Add a user permission to a space (admin only)",
		OperationID: "admin.addSpaceUserPermission",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/spaces/:space_id/permissions/users/:user_id", controller.AdminRemoveSpaceUserPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Remove user permission",
		Description: "Remove a user permission from a space (admin only)",
		OperationID: "admin.removeSpaceUserPermission",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Post(controller.FiberOapi, "/spaces/:space_id/permissions/groups", controller.AdminAddSpaceGroupPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Add group permission",
		Description: "Add a group permission to a space (admin only)",
		OperationID: "admin.addSpaceGroupPermission",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/spaces/:space_id/permissions/groups/:group_id", controller.AdminRemoveSpaceGroupPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Remove group permission",
		Description: "Remove a group permission from a space (admin only)",
		OperationID: "admin.removeSpaceGroupPermission",
		Tags:        []string{"Admin"},
	})

	// API Keys management
	fiberoapi.Get(controller.FiberOapi, "/apikeys", controller.ListAllApiKeys, fiberoapi.OpenAPIOptions{
		Summary:     "List all API keys",
		Description: "Retrieve a paginated list of all API keys (admin only)",
		OperationID: "admin.listAllApiKeys",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/apikeys/:apikey_id", controller.RevokeApiKey, fiberoapi.OpenAPIOptions{
		Summary:     "Revoke API key",
		Description: "Revoke an API key (admin only)",
		OperationID: "admin.revokeApiKey",
		Tags:        []string{"Admin"},
	})

	// Groups management
	fiberoapi.Get(controller.FiberOapi, "/groups", controller.ListGroups, fiberoapi.OpenAPIOptions{
		Summary:     "List all groups",
		Description: "Retrieve a paginated list of all groups (admin only)",
		OperationID: "admin.listGroups",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Post(controller.FiberOapi, "/groups", controller.CreateGroup, fiberoapi.OpenAPIOptions{
		Summary:     "Create group",
		Description: "Create a new group (admin only)",
		OperationID: "admin.createGroup",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Put(controller.FiberOapi, "/groups/:group_id", controller.UpdateGroup, fiberoapi.OpenAPIOptions{
		Summary:     "Update group",
		Description: "Update a group's name, description, or role (admin only)",
		OperationID: "admin.updateGroup",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/groups/:group_id", controller.DeleteGroup, fiberoapi.OpenAPIOptions{
		Summary:     "Delete group",
		Description: "Delete a group (admin only)",
		OperationID: "admin.deleteGroup",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Get(controller.FiberOapi, "/groups/:group_id/members", controller.GetGroupMembers, fiberoapi.OpenAPIOptions{
		Summary:     "Get group members",
		Description: "Retrieve all members of a group (admin only)",
		OperationID: "admin.getGroupMembers",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Post(controller.FiberOapi, "/groups/:group_id/members", controller.AddGroupMember, fiberoapi.OpenAPIOptions{
		Summary:     "Add group member",
		Description: "Add a user to a group (admin only)",
		OperationID: "admin.addGroupMember",
		Tags:        []string{"Admin"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/groups/:group_id/members/:user_id", controller.RemoveGroupMember, fiberoapi.OpenAPIOptions{
		Summary:     "Remove group member",
		Description: "Remove a user from a group (admin only)",
		OperationID: "admin.removeGroupMember",
		Tags:        []string{"Admin"},
	})
}
