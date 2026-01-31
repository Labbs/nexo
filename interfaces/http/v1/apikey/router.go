package apikey

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/apikey"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config    config.Config
	Logger    zerolog.Logger
	FiberOapi *fiberoapi.OApiGroup
	ApiKeyApp *apikey.ApiKeyApp
}

func SetupApiKeyRouter(controller Controller) {
	fiberoapi.Get(controller.FiberOapi, "/", controller.ListApiKeys, fiberoapi.OpenAPIOptions{
		Summary:     "List API keys",
		Description: "List all API keys for the current user",
		OperationID: "apikey.list",
		Tags:        []string{"API Keys"},
	})
	fiberoapi.Post(controller.FiberOapi, "/", controller.CreateApiKey, fiberoapi.OpenAPIOptions{
		Summary:     "Create API key",
		Description: "Create a new API key. The key is only shown once.",
		OperationID: "apikey.create",
		Tags:        []string{"API Keys"},
	})
	fiberoapi.Put(controller.FiberOapi, "/:api_key_id", controller.UpdateApiKey, fiberoapi.OpenAPIOptions{
		Summary:     "Update API key",
		Description: "Update an API key's name or scopes",
		OperationID: "apikey.update",
		Tags:        []string{"API Keys"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/:api_key_id", controller.DeleteApiKey, fiberoapi.OpenAPIOptions{
		Summary:     "Delete API key",
		Description: "Revoke and delete an API key",
		OperationID: "apikey.delete",
		Tags:        []string{"API Keys"},
	})
	fiberoapi.Get(controller.FiberOapi, "/scopes", controller.GetAvailableScopes, fiberoapi.OpenAPIOptions{
		Summary:     "Get available scopes",
		Description: "List all available permission scopes for API keys",
		OperationID: "apikey.scopes",
		Tags:        []string{"API Keys"},
	})
}
