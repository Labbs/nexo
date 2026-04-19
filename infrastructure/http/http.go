package http

import (
	"encoding/json"
	"time"

	"github.com/labbs/nexo/application/session"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/labbs/nexo/infrastructure/logger/zerolog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberoapi "github.com/labbs/fiber-oapi"
	z "github.com/rs/zerolog"
)

type Config struct {
	Fiber     *fiber.App
	FiberOapi *fiberoapi.OApiApp
}

// Configure sets up the HTTP server (fiber) with the provided configuration and logger.
// The FiberOapi instance is also configured for OpenAPI support and exposed documentation.
// Will return an error if the server cannot be created (fatal)
func Configure(_cfg config.Config, logger z.Logger, sessionApp *session.SessionApp, enableIU bool) (Config, error) {
	var c Config
	fiberConfig := fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		BodyLimit:             10 * 1024 * 1024, // 10 MB
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           120 * time.Second,
	}

	r := fiber.New(fiberConfig)

	if _cfg.Server.HttpLogs {
		r.Use(zerolog.HTTPLogger(logger))
	}

	r.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     _cfg.Server.CorsAllowOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
		MaxAge:           86400,
	}))
	r.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		return c.Next()
	})
	r.Use(compress.New())
	r.Use(requestid.New())

	authAdapter := NewSessionAuthAdapter(sessionApp)

	oapiConfig := fiberoapi.Config{
		EnableValidation:    true,
		EnableOpenAPIDocs:   true,
		OpenAPIDocsPath:     "/documentation",
		OpenAPIJSONPath:     "/api-spec.json",
		OpenAPIYamlPath:     "/api-spec.yaml",
		AuthService:         authAdapter,
		EnableAuthorization: true,
		SecuritySchemes: map[string]fiberoapi.SecurityScheme{
			"bearerAuth": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT Bearer token",
			},
		},
		DefaultSecurity: []map[string][]string{
			{"bearerAuth": {}},
		},
	}

	c.FiberOapi = fiberoapi.New(r, oapiConfig)
	c.Fiber = r

	return c, nil
}
