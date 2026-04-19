package collaboration

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/session"
	sessionDto "github.com/labbs/nexo/application/session/dto"
	"github.com/labbs/nexo/domain"
	"github.com/rs/zerolog"
)

const wsCollabPrefix = "/ws/collab/"

// Handler manages WebSocket connections for Y.js collaboration.
type Handler struct {
	hub          *Hub
	sessionApp   *session.SessionApplication
	documentPers domain.DocumentPers
	logger       zerolog.Logger
}

// NewHandler creates a new collaboration WebSocket handler.
func NewHandler(hub *Hub, sessionApp *session.SessionApplication, documentPers domain.DocumentPers, logger zerolog.Logger) *Handler {
	return &Handler{
		hub:          hub,
		sessionApp:   sessionApp,
		documentPers: documentPers,
		logger:       logger.With().Str("component", "collaboration.handler").Logger(),
	}
}

// UpgradeMiddleware validates the JWT token before upgrading the WebSocket connection.
func (h *Handler) UpgradeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		h.logger.Debug().Str("event", "ws_upgrade").Str("path", c.Path()).Msg("upgrading to WebSocket")

		token := c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		result, err := h.sessionApp.ValidateToken(sessionDto.ValidateTokenInput{Token: token})
		if err != nil {
			h.logger.Warn().Err(err).Msg("invalid token on websocket upgrade")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("auth_context", result.AuthContext)
		c.Locals("user_id", result.AuthContext.UserID)
		c.Locals("path", c.Path())

		return c.Next()
	}
}

// canAccessRoom checks if the user has at least viewer access to the room's resource.
// roomID format: "document:{id}", "drawing:{id}", "row:{databaseId}:{rowId}"
func (h *Handler) canAccessRoom(userID string, authCtx *fiberoapi.AuthContext, roomID string) bool {
	// Admin bypasses all checks
	for _, role := range authCtx.Roles {
		if role == string(domain.RoleAdmin) {
			return true
		}
	}

	parts := strings.SplitN(roomID, ":", 3)
	if len(parts) < 2 {
		return false
	}

	resourceType := parts[0]
	resourceID := parts[1]

	switch resourceType {
	case "document":
		doc, err := h.documentPers.GetDocumentWithPermissions(resourceID, userID)
		if err != nil {
			h.logger.Warn().Err(err).Str("room_id", roomID).Msg("failed to load document for access check")
			return false
		}
		return doc.HasPermission(userID, domain.PermissionRoleViewer)

	case "drawing":
		ok, err := h.sessionApp.CanAccessResource(sessionDto.CanAccessResourceInput{
			Context:      authCtx,
			ResourceType: "drawing",
			ResourceID:   resourceID,
			Action:       "read",
		})
		if err != nil {
			h.logger.Warn().Err(err).Str("room_id", roomID).Msg("failed to check drawing access")
			return false
		}
		return ok

	case "row", "database":
		ok, err := h.sessionApp.CanAccessResource(sessionDto.CanAccessResourceInput{
			Context:      authCtx,
			ResourceType: "database",
			ResourceID:   resourceID,
			Action:       "read",
		})
		if err != nil {
			h.logger.Warn().Err(err).Str("room_id", roomID).Msg("failed to check database access")
			return false
		}
		return ok

	default:
		return false
	}
}

// WebSocketHandler returns the Fiber WebSocket handler for collaboration.
func (h *Handler) WebSocketHandler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		roomID := strings.TrimPrefix(c.Locals("path").(string), wsCollabPrefix)
		userID, _ := c.Locals("user_id").(string)
		authCtx, _ := c.Locals("auth_context").(*fiberoapi.AuthContext)

		h.logger.Debug().Str("event", "ws_connection").Str("room_id", roomID).Str("user_id", userID).Msg("new WebSocket connection")

		if roomID == "" {
			h.logger.Warn().Msg("empty room id")
			return
		}

		if authCtx == nil || !h.canAccessRoom(userID, authCtx, roomID) {
			h.logger.Warn().Str("room_id", roomID).Str("user_id", userID).Msg("unauthorized WebSocket room access")
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "forbidden"))
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

		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					h.logger.Warn().Err(err).Str("room_id", roomID).Str("user_id", userID).Msg("unexpected close")
				}
				break
			}

			if messageType == websocket.BinaryMessage {
				room.Broadcast(c, msg)
			}
		}
	})
}
