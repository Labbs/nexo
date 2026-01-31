package database

import fiberoapi "github.com/labbs/fiber-oapi"

func SetupDatabaseRouter(ctrl Controller) {
	// Database endpoints
	fiberoapi.Get(ctrl.FiberOapi, "/", ctrl.ListDatabases, fiberoapi.OpenAPIOptions{
		Summary:     "List databases",
		Description: "List all databases in a space",
		OperationID: "database.list",
		Tags:        []string{"Databases"},
	})

	fiberoapi.Post(ctrl.FiberOapi, "/", ctrl.CreateDatabase, fiberoapi.OpenAPIOptions{
		Summary:     "Create database",
		Description: "Create a new database",
		OperationID: "database.create",
		Tags:        []string{"Databases"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/types", ctrl.GetAvailableTypes, fiberoapi.OpenAPIOptions{
		Summary:     "Get available property types",
		Description: "List all available property types for database columns",
		OperationID: "database.types",
		Tags:        []string{"Databases"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/search", ctrl.SearchDatabases, fiberoapi.OpenAPIOptions{
		Summary:     "Search databases",
		Description: "Search databases by name or description",
		OperationID: "database.search",
		Tags:        []string{"Databases", "Search"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/:database_id", ctrl.GetDatabase, fiberoapi.OpenAPIOptions{
		Summary:     "Get database",
		Description: "Get a specific database by ID",
		OperationID: "database.get",
		Tags:        []string{"Databases"},
	})

	fiberoapi.Put(ctrl.FiberOapi, "/:database_id", ctrl.UpdateDatabase, fiberoapi.OpenAPIOptions{
		Summary:     "Update database",
		Description: "Update an existing database",
		OperationID: "database.update",
		Tags:        []string{"Databases"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:database_id", ctrl.DeleteDatabase, fiberoapi.OpenAPIOptions{
		Summary:     "Delete database",
		Description: "Delete a database and all its rows",
		OperationID: "database.delete",
		Tags:        []string{"Databases"},
	})

	fiberoapi.Patch(ctrl.FiberOapi, "/:database_id/move", ctrl.MoveDatabase, fiberoapi.OpenAPIOptions{
		Summary:     "Move database",
		Description: "Move a database to a document or to root level",
		OperationID: "database.move",
		Tags:        []string{"Databases"},
	})

	// View endpoints
	fiberoapi.Post(ctrl.FiberOapi, "/:database_id/views", ctrl.CreateView, fiberoapi.OpenAPIOptions{
		Summary:     "Create view",
		Description: "Create a new view for a database",
		OperationID: "database.view.create",
		Tags:        []string{"Database Views"},
	})

	fiberoapi.Put(ctrl.FiberOapi, "/:database_id/views/:view_id", ctrl.UpdateView, fiberoapi.OpenAPIOptions{
		Summary:     "Update view",
		Description: "Update an existing view",
		OperationID: "database.view.update",
		Tags:        []string{"Database Views"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:database_id/views/:view_id", ctrl.DeleteView, fiberoapi.OpenAPIOptions{
		Summary:     "Delete view",
		Description: "Delete a view",
		OperationID: "database.view.delete",
		Tags:        []string{"Database Views"},
	})

	// Row endpoints
	fiberoapi.Get(ctrl.FiberOapi, "/:database_id/rows", ctrl.ListRows, fiberoapi.OpenAPIOptions{
		Summary:     "List rows",
		Description: "List all rows in a database",
		OperationID: "database.row.list",
		Tags:        []string{"Database Rows"},
	})

	fiberoapi.Post(ctrl.FiberOapi, "/:database_id/rows", ctrl.CreateRow, fiberoapi.OpenAPIOptions{
		Summary:     "Create row",
		Description: "Create a new row in a database",
		OperationID: "database.row.create",
		Tags:        []string{"Database Rows"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:database_id/rows", ctrl.BulkDeleteRows, fiberoapi.OpenAPIOptions{
		Summary:     "Bulk delete rows",
		Description: "Delete multiple rows at once",
		OperationID: "database.row.bulk_delete",
		Tags:        []string{"Database Rows"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/:database_id/rows/:row_id", ctrl.GetRow, fiberoapi.OpenAPIOptions{
		Summary:     "Get row",
		Description: "Get a specific row by ID",
		OperationID: "database.row.get",
		Tags:        []string{"Database Rows"},
	})

	fiberoapi.Put(ctrl.FiberOapi, "/:database_id/rows/:row_id", ctrl.UpdateRow, fiberoapi.OpenAPIOptions{
		Summary:     "Update row",
		Description: "Update an existing row",
		OperationID: "database.row.update",
		Tags:        []string{"Database Rows"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:database_id/rows/:row_id", ctrl.DeleteRow, fiberoapi.OpenAPIOptions{
		Summary:     "Delete row",
		Description: "Delete a row",
		OperationID: "database.row.delete",
		Tags:        []string{"Database Rows"},
	})

	// Permission endpoints
	fiberoapi.Get(ctrl.FiberOapi, "/:database_id/permissions", ctrl.ListDatabasePermissions, fiberoapi.OpenAPIOptions{
		Summary:     "List database permissions",
		Description: "List all permissions for a database",
		OperationID: "database.permission.list",
		Tags:        []string{"Database Permissions"},
	})

	fiberoapi.Post(ctrl.FiberOapi, "/:database_id/permissions", ctrl.UpsertDatabasePermission, fiberoapi.OpenAPIOptions{
		Summary:     "Add or update permission",
		Description: "Add or update a permission for a database",
		OperationID: "database.permission.upsert",
		Tags:        []string{"Database Permissions"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:database_id/permissions", ctrl.DeleteDatabasePermission, fiberoapi.OpenAPIOptions{
		Summary:     "Delete permission",
		Description: "Delete a permission from a database",
		OperationID: "database.permission.delete",
		Tags:        []string{"Database Permissions"},
	})
}
