package collaboration

import (
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

// Room represents a Y.js collaboration room.
// It acts as a pure relay: binary messages from one client are broadcast to all others.
type Room struct {
	id      string
	mu      sync.RWMutex
	clients map[*websocket.Conn]*Client
	logger  zerolog.Logger
}

// Client holds metadata about a connected user.
type Client struct {
	UserID   string
	Username string
}

func newRoom(id string, logger zerolog.Logger) *Room {
	return &Room{
		id:      id,
		clients: make(map[*websocket.Conn]*Client),
		logger:  logger.With().Str("component", "collaboration.room").Str("room_id", id).Logger(),
	}
}

// AddClient registers a new WebSocket connection in the room.
func (r *Room) AddClient(conn *websocket.Conn, client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[conn] = client
	r.logger.Info().Str("user_id", client.UserID).Int("clients", len(r.clients)).Msg("client joined")
}

// RemoveClient unregisters a WebSocket connection from the room.
func (r *Room) RemoveClient(conn *websocket.Conn) {
	r.mu.Lock()
	client, ok := r.clients[conn]
	if ok {
		delete(r.clients, conn)
	}
	count := len(r.clients)
	r.mu.Unlock()

	if ok {
		r.logger.Info().Str("user_id", client.UserID).Int("clients", count).Msg("client left")
	}
}

// Broadcast sends a binary message to all clients except the sender.
func (r *Room) Broadcast(sender *websocket.Conn, msg []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for conn := range r.clients {
		if conn == sender {
			continue
		}
		if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			r.logger.Warn().Err(err).Msg("failed to write to client")
		}
	}
}

// ClientCount returns the number of connected clients.
func (r *Room) ClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}
