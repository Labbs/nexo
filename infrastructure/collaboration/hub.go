package collaboration

import (
	"sync"

	"github.com/rs/zerolog"
)

// Hub manages all collaboration rooms.
type Hub struct {
	mu     sync.RWMutex
	rooms  map[string]*Room
	logger zerolog.Logger
}

// NewHub creates a new collaboration hub.
func NewHub(logger zerolog.Logger) *Hub {
	return &Hub{
		rooms:  make(map[string]*Room),
		logger: logger.With().Str("component", "collaboration.hub").Logger(),
	}
}

// GetOrCreateRoom returns an existing room or creates a new one.
func (h *Hub) GetOrCreateRoom(roomID string) *Room {
	h.mu.RLock()
	room, ok := h.rooms[roomID]
	h.mu.RUnlock()
	if ok {
		return room
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Double-check after acquiring write lock
	if room, ok = h.rooms[roomID]; ok {
		return room
	}

	room = newRoom(roomID, h.logger)
	h.rooms[roomID] = room
	h.logger.Info().Str("room_id", roomID).Msg("room created")
	return room
}

// RemoveRoomIfEmpty removes a room if it has no more clients.
func (h *Hub) RemoveRoomIfEmpty(roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.rooms[roomID]
	if !ok {
		return
	}

	if room.ClientCount() == 0 {
		delete(h.rooms, roomID)
		h.logger.Info().Str("room_id", roomID).Msg("room removed (empty)")
	}
}

// Stats returns the number of active rooms and total clients.
func (h *Hub) Stats() (rooms int, clients int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rooms = len(h.rooms)
	for _, r := range h.rooms {
		clients += r.ClientCount()
	}
	return
}
