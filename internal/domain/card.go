package domain

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID              uuid.UUID  `json:"id"`
	ProjectID       uuid.UUID  `json:"project_id"`
	BoardID         uuid.UUID  `json:"board_id"`
	ColumnID        uuid.UUID  `json:"column_id"`
	CardNumber      int        `json:"card_number"`
	Title           string     `json:"title"`
	Description     string     `json:"description,omitempty"`
	DescriptionJSON any        `json:"description_json,omitempty"`
	Priority        Priority   `json:"priority"`
	Position        string     `json:"position"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	DueDate         *time.Time `json:"due_date,omitempty"`
	ArchivedAt      *time.Time `json:"archived_at,omitempty"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Joined fields (not stored directly)
	Assignees []CardAssignee `json:"assignees,omitempty"`
	Labels    []Label        `json:"labels,omitempty"`

	// ReadableID is project_key + "-" + card_number (e.g. "PROJ-123")
	ReadableID string `json:"readable_id,omitempty"`
}

type CardAssignee struct {
	CardID    uuid.UUID `json:"card_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	// Joined
	UserName string `json:"user_name,omitempty"`
}
