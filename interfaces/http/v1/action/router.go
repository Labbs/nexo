package action

import fiberoapi "github.com/labbs/fiber-oapi"

func SetupActionRouter(ctrl Controller) {
	fiberoapi.Get(ctrl.FiberOapi, "/", ctrl.ListActions, fiberoapi.OpenAPIOptions{
		Summary:     "List actions",
		Description: "List all automation actions for the authenticated user",
		OperationID: "action.list",
		Tags:        []string{"Actions"},
	})

	fiberoapi.Post(ctrl.FiberOapi, "/", ctrl.CreateAction, fiberoapi.OpenAPIOptions{
		Summary:     "Create action",
		Description: "Create a new automation action",
		OperationID: "action.create",
		Tags:        []string{"Actions"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/triggers", ctrl.GetAvailableTriggers, fiberoapi.OpenAPIOptions{
		Summary:     "Get available triggers",
		Description: "List all available trigger types for actions",
		OperationID: "action.triggers",
		Tags:        []string{"Actions"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/steps", ctrl.GetAvailableSteps, fiberoapi.OpenAPIOptions{
		Summary:     "Get available steps",
		Description: "List all available step types for actions",
		OperationID: "action.steps",
		Tags:        []string{"Actions"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/:action_id", ctrl.GetAction, fiberoapi.OpenAPIOptions{
		Summary:     "Get action",
		Description: "Get a specific action by ID",
		OperationID: "action.get",
		Tags:        []string{"Actions"},
	})

	fiberoapi.Put(ctrl.FiberOapi, "/:action_id", ctrl.UpdateAction, fiberoapi.OpenAPIOptions{
		Summary:     "Update action",
		Description: "Update an existing action",
		OperationID: "action.update",
		Tags:        []string{"Actions"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:action_id", ctrl.DeleteAction, fiberoapi.OpenAPIOptions{
		Summary:     "Delete action",
		Description: "Delete an action",
		OperationID: "action.delete",
		Tags:        []string{"Actions"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/:action_id/runs", ctrl.GetRuns, fiberoapi.OpenAPIOptions{
		Summary:     "Get action runs",
		Description: "Get execution history for an action",
		OperationID: "action.runs",
		Tags:        []string{"Actions"},
	})
}
