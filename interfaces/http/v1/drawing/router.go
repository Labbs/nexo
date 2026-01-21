package drawing

import fiberoapi "github.com/labbs/fiber-oapi"

func SetupDrawingRouter(ctrl Controller) {
	fiberoapi.Get(ctrl.FiberOapi, "/", ctrl.ListDrawings, fiberoapi.OpenAPIOptions{
		Summary:     "List drawings",
		Description: "List all drawings in a space",
		OperationID: "drawing.list",
		Tags:        []string{"Drawings"},
	})

	fiberoapi.Post(ctrl.FiberOapi, "/", ctrl.CreateDrawing, fiberoapi.OpenAPIOptions{
		Summary:     "Create drawing",
		Description: "Create a new Excalidraw drawing",
		OperationID: "drawing.create",
		Tags:        []string{"Drawings"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/:drawing_id", ctrl.GetDrawing, fiberoapi.OpenAPIOptions{
		Summary:     "Get drawing",
		Description: "Get a specific drawing by ID",
		OperationID: "drawing.get",
		Tags:        []string{"Drawings"},
	})

	fiberoapi.Put(ctrl.FiberOapi, "/:drawing_id", ctrl.UpdateDrawing, fiberoapi.OpenAPIOptions{
		Summary:     "Update drawing",
		Description: "Update an existing drawing",
		OperationID: "drawing.update",
		Tags:        []string{"Drawings"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:drawing_id", ctrl.DeleteDrawing, fiberoapi.OpenAPIOptions{
		Summary:     "Delete drawing",
		Description: "Delete a drawing",
		OperationID: "drawing.delete",
		Tags:        []string{"Drawings"},
	})

	// Permission routes
	fiberoapi.Get(ctrl.FiberOapi, "/:drawing_id/permissions", ctrl.ListDrawingPermissions, fiberoapi.OpenAPIOptions{
		Summary:     "List drawing permissions",
		Description: "List all permissions for a drawing",
		OperationID: "drawing.permissions.list",
		Tags:        []string{"Drawings"},
	})

	fiberoapi.Put(ctrl.FiberOapi, "/:drawing_id/permissions/user", ctrl.UpsertDrawingUserPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Upsert drawing user permission",
		Description: "Create or update a user permission for a drawing",
		OperationID: "drawing.permissions.upsert_user",
		Tags:        []string{"Drawings"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:drawing_id/permissions/user/:user_id", ctrl.DeleteDrawingUserPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Delete drawing user permission",
		Description: "Delete a user permission from a drawing",
		OperationID: "drawing.permissions.delete_user",
		Tags:        []string{"Drawings"},
	})
}
