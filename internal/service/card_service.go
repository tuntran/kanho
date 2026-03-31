package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/fracidx"
	"github.com/tungtran/kanho/internal/repository"
)

type CardService struct {
	cardRepo     *repository.CardRepo
	activityRepo *repository.ActivityRepo
	projectRepo  *repository.ProjectRepo
	boardRepo    *repository.BoardRepo
}

func NewCardService(
	cardRepo *repository.CardRepo,
	activityRepo *repository.ActivityRepo,
	projectRepo *repository.ProjectRepo,
	boardRepo *repository.BoardRepo,
) *CardService {
	return &CardService{
		cardRepo:     cardRepo,
		activityRepo: activityRepo,
		projectRepo:  projectRepo,
		boardRepo:    boardRepo,
	}
}

type CreateCardInput struct {
	BoardID     uuid.UUID
	ColumnID    uuid.UUID
	Title       string
	Description string
	Priority    domain.Priority
	StartDate   *time.Time
	DueDate     *time.Time
	CreatedBy   uuid.UUID
}

func (s *CardService) Create(ctx context.Context, in CreateCardInput) (*domain.Card, error) {
	board, err := s.boardRepo.GetByID(ctx, in.BoardID)
	if err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	project, err := s.projectRepo.GetByID(ctx, board.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	// Determine position at end of column
	positions, err := s.cardRepo.Positions(ctx, in.ColumnID)
	if err != nil {
		return nil, fmt.Errorf("get positions: %w", err)
	}

	var position string
	if len(positions) == 0 {
		position = fracidx.Start()
	} else {
		position = fracidx.End(positions[len(positions)-1])
	}

	if in.Priority == "" {
		in.Priority = domain.PriorityNone
	}

	now := time.Now()
	card := &domain.Card{
		ID:          uuid.New(),
		ProjectID:   board.ProjectID,
		BoardID:     in.BoardID,
		ColumnID:    in.ColumnID,
		Title:       in.Title,
		Description: in.Description,
		Priority:    in.Priority,
		Position:    position,
		StartDate:   in.StartDate,
		DueDate:     in.DueDate,
		CreatedBy:   in.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.cardRepo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("create card: %w", err)
	}

	card.ReadableID = project.Key + "-" + fmt.Sprint(card.CardNumber)

	// Log activity
	s.logActivity(ctx, project.WorkspaceID, in.BoardID, board.ProjectID, in.CreatedBy,
		"card.created", "card", card.ID, map[string]any{
			"title":       card.Title,
			"readable_id": card.ReadableID,
			"column_id":   in.ColumnID.String(),
		})

	return card, nil
}

func (s *CardService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Card, error) {
	return s.cardRepo.GetByID(ctx, id)
}

func (s *CardService) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]domain.Card, error) {
	return s.cardRepo.ListByBoard(ctx, boardID)
}

func (s *CardService) ListByColumn(ctx context.Context, columnID uuid.UUID) ([]domain.Card, error) {
	return s.cardRepo.ListByColumn(ctx, columnID)
}

type UpdateCardInput struct {
	Title       *string
	Description *string
	Priority    *domain.Priority
	StartDate   *time.Time
	DueDate     *time.Time
}

func (s *CardService) Update(ctx context.Context, id uuid.UUID, in UpdateCardInput, actorID uuid.UUID) (*domain.Card, error) {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	changes := map[string]any{}

	if in.Title != nil && *in.Title != card.Title {
		changes["title"] = map[string]string{"from": card.Title, "to": *in.Title}
		card.Title = *in.Title
	}
	if in.Description != nil && *in.Description != card.Description {
		changes["description"] = "updated"
		card.Description = *in.Description
	}
	if in.Priority != nil && *in.Priority != card.Priority {
		changes["priority"] = map[string]string{"from": string(card.Priority), "to": string(*in.Priority)}
		card.Priority = *in.Priority
	}
	if in.StartDate != nil {
		card.StartDate = in.StartDate
	}
	if in.DueDate != nil {
		card.DueDate = in.DueDate
	}

	if err := s.cardRepo.Update(ctx, card); err != nil {
		return nil, err
	}

	if len(changes) > 0 {
		board, _ := s.boardRepo.GetByID(ctx, card.BoardID)
		var wsID uuid.UUID
		if board != nil {
			project, _ := s.projectRepo.GetByID(ctx, board.ProjectID)
			if project != nil {
				wsID = project.WorkspaceID
			}
		}
		s.logActivity(ctx, wsID, card.BoardID, card.ProjectID, actorID,
			"card.updated", "card", card.ID, changes)
	}

	return card, nil
}

type MoveCardInput struct {
	ColumnID    uuid.UUID
	AfterCardID *uuid.UUID // nil = move to top
	BeforeCardID *uuid.UUID // nil = move to bottom
}

func (s *CardService) Move(ctx context.Context, cardID uuid.UUID, in MoveCardInput, actorID uuid.UUID) error {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return fmt.Errorf("get card: %w", err)
	}

	positions, err := s.cardRepo.Positions(ctx, in.ColumnID)
	if err != nil {
		return fmt.Errorf("get positions: %w", err)
	}

	var newPosition string
	switch {
	case in.AfterCardID == nil && in.BeforeCardID == nil:
		// Move to end
		if len(positions) == 0 {
			newPosition = fracidx.Start()
		} else {
			newPosition = fracidx.End(positions[len(positions)-1])
		}
	case in.AfterCardID != nil && in.BeforeCardID != nil:
		// Between two cards — find their positions
		afterPos := findPosition(positions, *in.AfterCardID, ctx, s.cardRepo)
		beforePos := findPosition(positions, *in.BeforeCardID, ctx, s.cardRepo)
		if afterPos != "" && beforePos != "" {
			newPosition = fracidx.GenerateBetween(afterPos, beforePos)
		} else {
			newPosition = fracidx.Start()
		}
	case in.AfterCardID != nil:
		// After a specific card (move after it, could be at end)
		afterPos := findPosition(positions, *in.AfterCardID, ctx, s.cardRepo)
		if afterPos != "" {
			// Find next position after afterPos
			nextPos := ""
			for _, p := range positions {
				if p > afterPos {
					nextPos = p
					break
				}
			}
			if nextPos == "" {
				newPosition = fracidx.End(afterPos)
			} else {
				newPosition = fracidx.GenerateBetween(afterPos, nextPos)
			}
		} else {
			newPosition = fracidx.Start()
		}
	case in.BeforeCardID != nil:
		// Before a specific card (move to top or before it)
		beforePos := findPosition(positions, *in.BeforeCardID, ctx, s.cardRepo)
		if beforePos != "" {
			// Find previous position
			prevPos := ""
			for i := len(positions) - 1; i >= 0; i-- {
				if positions[i] < beforePos {
					prevPos = positions[i]
					break
				}
			}
			if prevPos == "" {
				newPosition = fracidx.GenerateBetween("", beforePos)
			} else {
				newPosition = fracidx.GenerateBetween(prevPos, beforePos)
			}
		} else {
			newPosition = fracidx.Start()
		}
	}

	fromColumnID := card.ColumnID

	if err := s.cardRepo.Move(ctx, cardID, in.ColumnID, newPosition); err != nil {
		return fmt.Errorf("move card: %w", err)
	}

	// Log activity
	board, _ := s.boardRepo.GetByID(ctx, card.BoardID)
	var wsID uuid.UUID
	if board != nil {
		project, _ := s.projectRepo.GetByID(ctx, board.ProjectID)
		if project != nil {
			wsID = project.WorkspaceID
		}
	}
	s.logActivity(ctx, wsID, card.BoardID, card.ProjectID, actorID,
		"card.moved", "card", cardID, map[string]any{
			"readable_id":    card.ReadableID,
			"from_column_id": fromColumnID.String(),
			"to_column_id":   in.ColumnID.String(),
			"new_position":   newPosition,
		})

	return nil
}

func (s *CardService) Archive(ctx context.Context, id uuid.UUID, actorID uuid.UUID) error {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.cardRepo.Archive(ctx, id); err != nil {
		return err
	}

	board, _ := s.boardRepo.GetByID(ctx, card.BoardID)
	var wsID uuid.UUID
	if board != nil {
		project, _ := s.projectRepo.GetByID(ctx, board.ProjectID)
		if project != nil {
			wsID = project.WorkspaceID
		}
	}
	s.logActivity(ctx, wsID, card.BoardID, card.ProjectID, actorID,
		"card.archived", "card", id, map[string]any{"readable_id": card.ReadableID})
	return nil
}

func (s *CardService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.cardRepo.Delete(ctx, id)
}

func (s *CardService) AddAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
	return s.cardRepo.AddAssignee(ctx, cardID, userID)
}

func (s *CardService) RemoveAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
	return s.cardRepo.RemoveAssignee(ctx, cardID, userID)
}

func (s *CardService) AddLabel(ctx context.Context, cardID, labelID uuid.UUID) error {
	return s.cardRepo.AddLabel(ctx, cardID, labelID)
}

func (s *CardService) RemoveLabel(ctx context.Context, cardID, labelID uuid.UUID) error {
	return s.cardRepo.RemoveLabel(ctx, cardID, labelID)
}

func (s *CardService) ListActivities(ctx context.Context, cardID uuid.UUID, limit, offset int) ([]domain.Activity, error) {
	return s.activityRepo.ListByCard(ctx, cardID, limit, offset)
}

func (s *CardService) logActivity(ctx context.Context, wsID, boardID, projectID, actorID uuid.UUID, action, targetType string, targetID uuid.UUID, data map[string]any) {
	a := &domain.Activity{
		ID:          uuid.New(),
		WorkspaceID: wsID,
		ActorID:     actorID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		ProjectID:   projectID,
		BoardID:     boardID,
		Data:        data,
		CreatedAt:   time.Now(),
	}
	// Best-effort: don't fail the main operation if activity logging fails
	_ = s.activityRepo.Create(ctx, a)
}

// findPosition looks up the position of a card by ID from the positions slice.
// Falls back to DB lookup.
func findPosition(positions []string, cardID uuid.UUID, ctx context.Context, repo *repository.CardRepo) string {
	c, err := repo.GetByID(ctx, cardID)
	if err != nil {
		return ""
	}
	return c.Position
}
