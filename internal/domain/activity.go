package domain

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	ActorID     uuid.UUID `json:"actor_id"`
	Action      string    `json:"action"`
	TargetType  string    `json:"target_type"`
	TargetID    uuid.UUID `json:"target_id"`
	ProjectID   uuid.UUID `json:"project_id,omitempty"`
	BoardID     uuid.UUID `json:"board_id,omitempty"`
	Data        any       `json:"data"`
	CreatedAt   time.Time `json:"created_at"`

	// Joined
	ActorName string `json:"actor_name,omitempty"`
}
