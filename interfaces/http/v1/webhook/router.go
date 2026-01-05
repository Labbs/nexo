package webhook

import fiberoapi "github.com/labbs/fiber-oapi"

func SetupWebhookRouter(ctrl Controller) {
	fiberoapi.Get(ctrl.FiberOapi, "/", ctrl.ListWebhooks, fiberoapi.OpenAPIOptions{
		Summary:     "List webhooks",
		Description: "List all webhooks for the authenticated user",
		OperationID: "webhook.list",
		Tags:        []string{"Webhooks"},
	})

	fiberoapi.Post(ctrl.FiberOapi, "/", ctrl.CreateWebhook, fiberoapi.OpenAPIOptions{
		Summary:     "Create webhook",
		Description: "Create a new webhook",
		OperationID: "webhook.create",
		Tags:        []string{"Webhooks"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/events", ctrl.GetAvailableEvents, fiberoapi.OpenAPIOptions{
		Summary:     "Get available events",
		Description: "Get list of available webhook events",
		OperationID: "webhook.events",
		Tags:        []string{"Webhooks"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/:webhook_id", ctrl.GetWebhook, fiberoapi.OpenAPIOptions{
		Summary:     "Get webhook",
		Description: "Get a specific webhook by ID",
		OperationID: "webhook.get",
		Tags:        []string{"Webhooks"},
	})

	fiberoapi.Put(ctrl.FiberOapi, "/:webhook_id", ctrl.UpdateWebhook, fiberoapi.OpenAPIOptions{
		Summary:     "Update webhook",
		Description: "Update an existing webhook",
		OperationID: "webhook.update",
		Tags:        []string{"Webhooks"},
	})

	fiberoapi.Delete(ctrl.FiberOapi, "/:webhook_id", ctrl.DeleteWebhook, fiberoapi.OpenAPIOptions{
		Summary:     "Delete webhook",
		Description: "Delete a webhook",
		OperationID: "webhook.delete",
		Tags:        []string{"Webhooks"},
	})

	fiberoapi.Get(ctrl.FiberOapi, "/:webhook_id/deliveries", ctrl.GetDeliveries, fiberoapi.OpenAPIOptions{
		Summary:     "Get webhook deliveries",
		Description: "Get delivery history for a webhook",
		OperationID: "webhook.deliveries",
		Tags:        []string{"Webhooks"},
	})
}
