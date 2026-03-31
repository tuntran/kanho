package domain

import (
	"time"

	"github.com/google/uuid"
)

type Column struct {
	ID        uuid.UUID `json:"id"`
	BoardID   uuid.UUID `json:"board_id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	Color     string    `json:"color,omitempty"`
	WIPLimit  *int      `json:"wip_limit,omitempty"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
}
