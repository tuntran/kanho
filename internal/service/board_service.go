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

type BoardService struct {
	boardRepo  *repository.BoardRepo
	columnRepo *repository.ColumnRepo
}

func NewBoardService(boardRepo *repository.BoardRepo, columnRepo *repository.ColumnRepo) *BoardService {
	return &BoardService{boardRepo: boardRepo, columnRepo: columnRepo}
}

type CreateBoardInput struct {
	ProjectID   uuid.UUID
	Name        string
	Description string
	Visibility  domain.Visibility
}

func (s *BoardService) Create(ctx context.Context, in CreateBoardInput) (*domain.Board, error) {
	if in.Visibility == "" {
		in.Visibility = domain.VisibilityWorkspace
	}

	slug, err := domain.UniqueSlug(in.Name, func(slug string) (bool, error) {
		return s.boardRepo.SlugExists(ctx, in.ProjectID, slug)
	})
	if err != nil {
		return nil, fmt.Errorf("generate slug: %w", err)
	}

	now := time.Now()
	b := &domain.Board{
		ID:          uuid.New(),
		ProjectID:   in.ProjectID,
		Name:        in.Name,
		Slug:        slug,
		Description: in.Description,
		Visibility:  in.Visibility,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.boardRepo.Create(ctx, b); err != nil {
		return nil, fmt.Errorf("create board: %w", err)
	}

	// Auto-create 4 default columns
	if err := s.createDefaultColumns(ctx, b.ID, now); err != nil {
		return nil, fmt.Errorf("create default columns: %w", err)
	}

	return b, nil
}

func (s *BoardService) createDefaultColumns(ctx context.Context, boardID uuid.UUID, now time.Time) error {
	defaults := []struct {
		name   string
		isDone bool
	}{
		{"Backlog", false},
		{"To Do", false},
		{"In Progress", false},
		{"Done", true},
	}

	// Generate positions using fracidx
	positions := make([]string, len(defaults))
	positions[0] = fracidx.Start()
	for i := 1; i < len(defaults); i++ {
		positions[i] = fracidx.End(positions[i-1])
	}

	for i, d := range defaults {
		col := &domain.Column{
			ID:        uuid.New(),
			BoardID:   boardID,
			Name:      d.name,
			Position:  positions[i],
			IsDone:    d.isDone,
			CreatedAt: now,
		}
		if err := s.columnRepo.Create(ctx, col); err != nil {
			return err
		}
	}
	return nil
}

func (s *BoardService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Board, error) {
	return s.boardRepo.GetByID(ctx, id)
}

func (s *BoardService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.Board, error) {
	return s.boardRepo.ListByProject(ctx, projectID)
}

func (s *BoardService) Update(ctx context.Context, id uuid.UUID, name, description string) (*domain.Board, error) {
	b, err := s.boardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		b.Name = name
	}
	if description != "" {
		b.Description = description
	}
	if err := s.boardRepo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BoardService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.boardRepo.Delete(ctx, id)
}
