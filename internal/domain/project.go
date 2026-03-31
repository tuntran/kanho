package domain

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	CardCounter int       `json:"card_counter"`
	CreatedAt   time.Time `json:"created_at"`
}
