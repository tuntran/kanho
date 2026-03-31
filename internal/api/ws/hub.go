package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Hub manages WebSocket rooms keyed by board ID.
// A single goroutine listens on PostgreSQL "board_events" channel
// and fans out notifications to subscribed clients.
type Hub struct {
	rooms map[string]*Room
	mu    sync.RWMutex
	db    *pgxpool.Pool
}

// Room holds the set of clients subscribed to a specific board.
type Room struct {
	boardID string
	clients map[*Client]struct{}
	mu      sync.RWMutex
}

// NewHub creates a new WebSocket hub.
func NewHub(db *pgxpool.Pool) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		db:    db,
	}
}

// Start begins the LISTEN loop and cleanup goroutine. Call with `go hub.Start(ctx)`.
func (h *Hub) Start(ctx context.Context) {
	go h.listenLoop(ctx)
	go h.cleanupLoop(ctx)
}

// listenLoop acquires a dedicated connection and LISTENs on board_events.
func (h *Hub) listenLoop(ctx context.Context) {
	for {
		if err := h.listen(ctx); err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("ws hub listen error, reconnecting", "error", err)
			time.Sleep(2 * time.Second)
		}
	}
}

func (h *Hub) listen(ctx context.Context) error {
	conn, err := h.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "LISTEN board_events")
	if err != nil {
		return err
	}
	slog.Info("ws hub listening on board_events")

	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			return err
		}

		var payload struct {
			BoardID string `json:"board_id"`
		}
		if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
			slog.Warn("ws hub: invalid notification payload", "error", err)
			continue
		}

		if payload.BoardID == "" {
			continue
		}

		h.broadcast(payload.BoardID, []byte(notification.Payload))
	}
}

// broadcast sends data to all clients in the given board room.
func (h *Hub) broadcast(boardID string, data []byte) {
	h.mu.RLock()
	room, ok := h.rooms[boardID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	for client := range room.clients {
		select {
		case client.send <- data:
		default:
			// Client buffer full — drop and close
			go client.Close()
		}
	}
}

// Subscribe adds a client to a board room.
func (h *Hub) Subscribe(client *Client, boardID string) {
	h.mu.Lock()
	room, ok := h.rooms[boardID]
	if !ok {
		room = &Room{
			boardID: boardID,
			clients: make(map[*Client]struct{}),
		}
		h.rooms[boardID] = room
	}
	h.mu.Unlock()

	room.mu.Lock()
	room.clients[client] = struct{}{}
	room.mu.Unlock()

	client.mu.Lock()
	client.rooms[boardID] = struct{}{}
	client.mu.Unlock()
}

// Unsubscribe removes a client from a board room.
func (h *Hub) Unsubscribe(client *Client, boardID string) {
	h.mu.RLock()
	room, ok := h.rooms[boardID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	room.mu.Lock()
	delete(room.clients, client)
	empty := len(room.clients) == 0
	room.mu.Unlock()

	if empty {
		h.mu.Lock()
		// Double-check after acquiring write lock
		room.mu.RLock()
		if len(room.clients) == 0 {
			delete(h.rooms, boardID)
		}
		room.mu.RUnlock()
		h.mu.Unlock()
	}

	client.mu.Lock()
	delete(client.rooms, boardID)
	client.mu.Unlock()
}

// RemoveClient unsubscribes a client from all rooms.
func (h *Hub) RemoveClient(client *Client) {
	client.mu.RLock()
	boardIDs := make([]string, 0, len(client.rooms))
	for id := range client.rooms {
		boardIDs = append(boardIDs, id)
	}
	client.mu.RUnlock()

	for _, id := range boardIDs {
		h.Unsubscribe(client, id)
	}
}

// cleanupLoop periodically removes empty rooms.
func (h *Hub) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.mu.Lock()
			for id, room := range h.rooms {
				room.mu.RLock()
				if len(room.clients) == 0 {
					delete(h.rooms, id)
				}
				room.mu.RUnlock()
			}
			h.mu.Unlock()
		}
	}
}
