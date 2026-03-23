package collaboration

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/labbs/nexo/application/session"
	"github.com/labbs/nexo/application/session/dto"
	"github.com/rs/zerolog"
)

// Handler manages WebSocket connections for Y.js collaboration.
type Handler struct {
	hub        *Hub
	sessionApp *session.SessionApplication
	logger     zerolog.Logger
}

// NewHandler creates a new collaboration WebSocket handler.
func NewHandler(hub *Hub, sessionApp *session.SessionApplication, logger zerolog.Logger) *Handler {
	return &Handler{
		hub:        hub,
		sessionApp: sessionApp,
		logger:     logger.With().Str("component", "collaboration.handler").Logger(),
	}
}

// UpgradeMiddleware checks for WebSocket upgrade and validates the JWT token
// before upgrading the connection.
func (h *Handler) UpgradeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		token := c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		result, err := h.sessionApp.ValidateToken(dto.ValidateTokenInput{Token: token})
		if err != nil {
			h.logger.Warn().Err(err).Msg("invalid token on websocket upgrade")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		// Store auth context in locals for the WebSocket handler
		c.Locals("user_id", result.AuthContext.UserID)

		return c.Next()
	}
}

// WebSocketHandler returns the Fiber WebSocket handler for collaboration.
func (h *Handler) WebSocketHandler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		roomID := c.Params("roomId")
		userID, _ := c.Locals("user_id").(string)

		if roomID == "" {
			h.logger.Warn().Msg("empty room id")
			return
		}

		room := h.hub.GetOrCreateRoom(roomID)
		client := &Client{
			UserID: userID,
		}

		room.AddClient(c, client)
		defer func() {
			room.RemoveClient(c)
			h.hub.RemoveRoomIfEmpty(roomID)
		}()

		// Read loop: relay all binary messages to other clients in the room
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					h.logger.Warn().Err(err).Str("room_id", roomID).Str("user_id", userID).Msg("unexpected close")
				}
				break
			}

			// Only relay binary messages (Y.js protocol)
			if messageType == websocket.BinaryMessage {
				room.Broadcast(c, msg)
			}
		}
	})
}
