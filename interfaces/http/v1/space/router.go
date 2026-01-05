package space

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/space"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config    config.Config
	Logger    zerolog.Logger
	FiberOapi *fiberoapi.OApiGroup
	SpaceApp  *space.SpaceApp
}

func SetupSpaceRouter(controller Controller) {
	fiberoapi.Post(controller.FiberOapi, "", controller.CreateSpace, fiberoapi.OpenAPIOptions{
		Summary:     "Create a new space",
		Description: "Create a new space for the authenticated user",
		OperationID: "space.createSpace",
		Tags:        []string{"Space"},
	})

	fiberoapi.Put(controller.FiberOapi, "/:space_id", controller.UpdateSpace, fiberoapi.OpenAPIOptions{
		Summary:     "Update a space",
		Description: "Update space properties",
		OperationID: "space.updateSpace",
		Tags:        []string{"Space"},
	})

	fiberoapi.Delete(controller.FiberOapi, "/:space_id", controller.DeleteSpace, fiberoapi.OpenAPIOptions{
		Summary:     "Delete a space",
		Description: "Soft delete a space",
		OperationID: "space.deleteSpace",
		Tags:        []string{"Space"},
	})

	// Permissions (MVP: user-level)
	fiberoapi.Get(controller.FiberOapi, "/:space_id/permissions", controller.ListPermissions, fiberoapi.OpenAPIOptions{
		Summary:     "List space permissions",
		Description: "List user permissions for a space",
		OperationID: "space.listPermissions",
		Tags:        []string{"Space"},
	})

	fiberoapi.Put(controller.FiberOapi, "/:space_id/permissions", controller.UpsertUserPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Upsert user permission on space",
		Description: "Create or update a user's permission on a space",
		OperationID: "space.upsertUserPermission",
		Tags:        []string{"Space"},
	})

	fiberoapi.Delete(controller.FiberOapi, "/:space_id/permissions/:user_id", controller.DeleteUserPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Delete user permission on space",
		Description: "Remove user's permission from a space",
		OperationID: "space.deleteUserPermission",
		Tags:        []string{"Space"},
	})
}
