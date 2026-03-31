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

type ColumnService struct {
	repo *repository.ColumnRepo
}

func NewColumnService(repo *repository.ColumnRepo) *ColumnService {
	return &ColumnService{repo: repo}
}

type CreateColumnInput struct {
	BoardID  uuid.UUID
	Name     string
	Color    string
	WIPLimit *int
	IsDone   bool
}

func (s *ColumnService) Create(ctx context.Context, in CreateColumnInput) (*domain.Column, error) {
	// Get existing columns to determine position
	existing, err := s.repo.ListByBoard(ctx, in.BoardID)
	if err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}

	var position string
	if len(existing) == 0 {
		position = fracidx.Start()
	} else {
		lastPos := existing[len(existing)-1].Position
		position = fracidx.End(lastPos)
	}

	c := &domain.Column{
		ID:        uuid.New(),
		BoardID:   in.BoardID,
		Name:      in.Name,
		Position:  position,
		Color:     in.Color,
		WIPLimit:  in.WIPLimit,
		IsDone:    in.IsDone,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create column: %w", err)
	}
	return c, nil
}

func (s *ColumnService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Column, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ColumnService) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]domain.Column, error) {
	return s.repo.ListByBoard(ctx, boardID)
}

type UpdateColumnInput struct {
	Name     string
	Color    string
	WIPLimit *int
	IsDone   *bool
}

func (s *ColumnService) Update(ctx context.Context, id uuid.UUID, in UpdateColumnInput) (*domain.Column, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Name != "" {
		c.Name = in.Name
	}
	if in.Color != "" {
		c.Color = in.Color
	}
	if in.WIPLimit != nil {
		c.WIPLimit = in.WIPLimit
	}
	if in.IsDone != nil {
		c.IsDone = *in.IsDone
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ColumnService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

type ReorderItem struct {
	ID       uuid.UUID `json:"id"`
	Position string    `json:"position"`
}

func (s *ColumnService) Reorder(ctx context.Context, items []ReorderItem) error {
	if len(items) == 0 {
		return nil
	}

	// Validate positions are lexicographically ordered
	for i := 1; i < len(items); i++ {
		if items[i].Position <= items[i-1].Position {
			return fmt.Errorf("positions must be lexicographically ordered")
		}
	}

	updates := make([]domain.Column, len(items))
	for i, item := range items {
		updates[i] = domain.Column{ID: item.ID, Position: item.Position}
	}
	return s.repo.BatchUpdatePositions(ctx, updates)
}
