package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/repository"
)

type LabelService struct {
	repo *repository.LabelRepo
}

func NewLabelService(repo *repository.LabelRepo) *LabelService {
	return &LabelService{repo: repo}
}

type CreateLabelInput struct {
	WorkspaceID uuid.UUID
	Name        string
	Color       string
}

func (s *LabelService) Create(ctx context.Context, in CreateLabelInput) (*domain.Label, error) {
	if in.Color == "" {
		in.Color = "#6366f1"
	}

	l := &domain.Label{
		ID:          uuid.New(),
		WorkspaceID: in.WorkspaceID,
		Name:        in.Name,
		Color:       in.Color,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, l); err != nil {
		return nil, fmt.Errorf("create label: %w", err)
	}
	return l, nil
}

func (s *LabelService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Label, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *LabelService) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.Label, error) {
	return s.repo.ListByWorkspace(ctx, workspaceID)
}

type UpdateLabelInput struct {
	Name  string
	Color string
}

func (s *LabelService) Update(ctx context.Context, id uuid.UUID, in UpdateLabelInput) (*domain.Label, error) {
	l, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Name != "" {
		l.Name = in.Name
	}
	if in.Color != "" {
		l.Color = in.Color
	}

	if err := s.repo.Update(ctx, l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *LabelService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
