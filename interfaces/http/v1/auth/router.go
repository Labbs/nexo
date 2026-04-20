package auth

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/auth"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config          config.Config
	Logger          zerolog.Logger
	FiberOapi       *fiberoapi.OApiGroup
	AuthApplication *auth.AuthApplication
}

func SetupAuthRouter(controller Controller) {
	fiberoapi.Post(controller.FiberOapi, "/login", controller.Login, fiberoapi.OpenAPIOptions{
		Summary:     "Login user",
		Description: "Authenticate user and return a token",
		OperationID: "auth.login",
		Tags:        []string{"Auth"},
		Security:    "disabled",
	})

	fiberoapi.Get(controller.FiberOapi, "/logout", controller.Logout, fiberoapi.OpenAPIOptions{
		Summary:     "Logout user",
		Description: "Invalidate user session",
		OperationID: "auth.logout",
		Tags:        []string{"Auth"},
	})

	fiberoapi.Post(controller.FiberOapi, "/register", controller.Register, fiberoapi.OpenAPIOptions{
		Summary:     "Register user",
		Description: "Register a new user",
		OperationID: "auth.register",
		Tags:        []string{"Auth"},
		Security:    "disabled",
	})

	fiberoapi.Get(controller.FiberOapi, "/sso/redirect", controller.SSORedirect, fiberoapi.OpenAPIOptions{
		Summary:     "SSO redirect URL",
		Description: "Returns the provider authorization URL for SSO login",
		OperationID: "auth.sso.redirect",
		Tags:        []string{"Auth"},
		Security:    "disabled",
	})

	fiberoapi.Post(controller.FiberOapi, "/sso/callback", controller.SSOCallback, fiberoapi.OpenAPIOptions{
		Summary:     "SSO callback",
		Description: "Exchange OAuth2 code for a Nexo session token",
		OperationID: "auth.sso.callback",
		Tags:        []string{"Auth"},
		Security:    "disabled",
	})
}
