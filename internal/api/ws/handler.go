package ws

import (
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
	"github.com/tungtran/kanho/internal/auth"
)

// Handler handles WebSocket upgrade requests.
type Handler struct {
	hub       *Hub
	jwtSecret string
}

// NewHandler creates a new WebSocket handler.
func NewHandler(hub *Hub, jwtSecret string) *Handler {
	return &Handler{hub: hub, jwtSecret: jwtSecret}
}

// ServeHTTP upgrades the HTTP connection to WebSocket.
// JWT token is passed via query parameter: ?token=<jwt>
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(h.jwtSecret, token)
	if err != nil {
		http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow all origins in dev; tighten in production
	})
	if err != nil {
		slog.Error("ws upgrade failed", "error", err)
		return
	}

	client := NewClient(h.hub, conn, claims.UserID)
	slog.Info("ws client connected", "user_id", claims.UserID)

	// Run blocks until disconnect
	client.Run(r.Context())
}
