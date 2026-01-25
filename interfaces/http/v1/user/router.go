package user

import (
	fiberoapi "github.com/labbs/fiber-oapi"
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
}

func SetupUserRouter(controller Controller) {
	fiberoapi.Get(controller.FiberOapi, "/profile", controller.GetProfile, fiberoapi.OpenAPIOptions{
		Summary:     "Get user profile",
		Description: "Retrieve the profile of the authenticated user",
		OperationID: "user.getProfile",
		Tags:        []string{"User"},
	})
	fiberoapi.Get(controller.FiberOapi, "/my-spaces", controller.GetMySpaces, fiberoapi.OpenAPIOptions{
		Summary:     "Get my spaces",
		Description: "Retrieve the spaces of the authenticated user",
		OperationID: "user.getMySpaces",
		Tags:        []string{"User"},
	})
	fiberoapi.Get(controller.FiberOapi, "/my-favorites", controller.GetMyFavorites, fiberoapi.OpenAPIOptions{
		Summary:     "Get my favorites",
		Description: "Retrieve the favorite items of the authenticated user",
		OperationID: "user.getMyFavorites",
		Tags:        []string{"User"},
	})
	fiberoapi.Post(controller.FiberOapi, "/favorite/:space_id/:document_id", controller.AddFavorite, fiberoapi.OpenAPIOptions{
		Summary:     "Add favorite",
		Description: "Add a document to the user's favorites",
		OperationID: "user.addFavorite",
		Tags:        []string{"User"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/favorite/:favorite_id", controller.RemoveFavorite, fiberoapi.OpenAPIOptions{
		Summary:     "Remove favorite",
		Description: "Remove a document from the user's favorites",
		OperationID: "user.removeFavorite",
		Tags:        []string{"User"},
	})

	fiberoapi.Put(controller.FiberOapi, "/favorite/:favorite_id/position", controller.UpdateFavoritePosition, fiberoapi.OpenAPIOptions{
		Summary:     "Update favorite position",
		Description: "Reorder favorites by updating a favorite's position",
		OperationID: "user.updateFavoritePosition",
		Tags:        []string{"User"},
	})

	// Profile management
	fiberoapi.Put(controller.FiberOapi, "/profile", controller.UpdateProfile, fiberoapi.OpenAPIOptions{
		Summary:     "Update user profile",
		Description: "Update the authenticated user's profile (username, avatar, preferences)",
		OperationID: "user.updateProfile",
		Tags:        []string{"User"},
	})
	fiberoapi.Post(controller.FiberOapi, "/change-password", controller.ChangePassword, fiberoapi.OpenAPIOptions{
		Summary:     "Change password",
		Description: "Change the authenticated user's password",
		OperationID: "user.changePassword",
		Tags:        []string{"User"},
	})

	fiberoapi.Get(controller.FiberOapi, "/list", controller.ListUsers, fiberoapi.OpenAPIOptions{
		Summary:     "List users",
		Description: "Get a simplified list of all users (id, username, avatar) for use in person pickers",
		OperationID: "user.listUsers",
		Tags:        []string{"User"},
	})
}
