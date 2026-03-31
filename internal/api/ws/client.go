package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const (
	sendBufferSize = 256
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingInterval   = 50 * time.Second
)

// Client represents a single WebSocket connection.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	userID uuid.UUID
	rooms  map[string]struct{}
	send   chan []byte
	done   chan struct{}
	once   sync.Once
	mu     sync.RWMutex
}

// NewClient creates a client for the given WebSocket connection.
func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		userID: userID,
		rooms:  make(map[string]struct{}),
		send:   make(chan []byte, sendBufferSize),
		done:   make(chan struct{}),
	}
}

// Run starts read and write pumps. Blocks until client disconnects.
func (c *Client) Run(ctx context.Context) {
	go c.writePump(ctx)
	c.readPump(ctx)
}

// Close gracefully shuts down the client.
func (c *Client) Close() {
	c.once.Do(func() {
		close(c.done)
		c.hub.RemoveClient(c)
		c.conn.Close(websocket.StatusNormalClosure, "bye")
	})
}

// clientMessage is the message format clients send to subscribe/unsubscribe.
type clientMessage struct {
	Action  string `json:"action"`
	BoardID string `json:"board_id"`
}

// readPump reads messages from the WebSocket connection.
func (c *Client) readPump(ctx context.Context) {
	defer c.Close()

	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			return
		}

		var msg clientMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		switch msg.Action {
		case "subscribe":
			if msg.BoardID != "" {
				c.hub.Subscribe(c, msg.BoardID)
				slog.Debug("ws client subscribed", "board_id", msg.BoardID, "user_id", c.userID)
			}
		case "unsubscribe":
			if msg.BoardID != "" {
				c.hub.Unsubscribe(c, msg.BoardID)
			}
		}
	}
}

// writePump writes messages from the send channel to the WebSocket connection.
func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case <-c.done:
			return
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Ping(pingCtx)
			cancel()
			if err != nil {
				return
			}
		}
	}
}
