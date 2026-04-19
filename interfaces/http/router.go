package http

import (
	"github.com/labbs/nexo/infrastructure"
	"github.com/labbs/nexo/infrastructure/collaboration"
	v1 "github.com/labbs/nexo/interfaces/http/v1"
)

func SetupRoutes(deps infrastructure.Deps) {
	logger := deps.Logger.With().Str("component", "http.router").Logger()
	logger.Info().Str("event", "setup_routes").Msg("Setting up HTTP routes")

	// Setup system routes (health, metrics, etc.)
	setupSystemRoutes(deps)

	// Setup v1 routes
	v1.SetupRouterV1(deps)

	// Setup WebSocket collaboration route
	setupCollaborationRoutes(deps)
}

func setupCollaborationRoutes(deps infrastructure.Deps) {
	logger := deps.Logger.With().Str("component", "http.router.collaboration").Logger()
	logger.Info().Str("event", "setup_collaboration_routes").Msg("Setting up collaboration WebSocket routes")
	handler := collaboration.NewHandler(deps.CollaborationHub, deps.SessionApplication, deps.Logger)

	// The frontend connects to ws://<host>/<roomId>?token=<jwt>
	// Room formats: "document:<docId>" or "row:<databaseId>:<rowId>"
	// Use("/ws/collab") is a prefix match, Get uses wildcard for the room ID (contains colons)
	deps.Http.Fiber.Use("/ws/collab", handler.UpgradeMiddleware())
	deps.Http.Fiber.Get("/ws/collab/+", handler.WebSocketHandler())

	logger.Debug().Interface("paths", deps.Http.Fiber.GetRoutes()).Msg("Registered HTTP routes")
}
