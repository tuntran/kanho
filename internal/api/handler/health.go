package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/storage"
)

type HealthHandler struct {
	db      *pgxpool.Pool
	storage storage.Client
}

func NewHealthHandler(db *pgxpool.Pool, storage storage.Client) *HealthHandler {
	return &HealthHandler{db: db, storage: storage}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result := map[string]string{
		"db":      "ok",
		"storage": "ok",
	}
	status := http.StatusOK

	if err := h.db.Ping(ctx); err != nil {
		result["db"] = "error"
		status = http.StatusServiceUnavailable
	}

	if err := h.storage.Ping(ctx); err != nil {
		result["storage"] = "error"
		status = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}
