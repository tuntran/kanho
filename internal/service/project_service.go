package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/repository"
)

var projectKeyRegex = regexp.MustCompile(`^[A-Z]{2,10}$`)

type ProjectService struct {
	repo     *repository.ProjectRepo
	boardSvc *BoardService
}

func NewProjectService(repo *repository.ProjectRepo, boardSvc *BoardService) *ProjectService {
	return &ProjectService{repo: repo, boardSvc: boardSvc}
}

type CreateProjectInput struct {
	WorkspaceID uuid.UUID
	Name        string
	Key         string
}

func (s *ProjectService) Create(ctx context.Context, in CreateProjectInput) (*domain.Project, error) {
	if !projectKeyRegex.MatchString(in.Key) {
		return nil, fmt.Errorf("project key must be 2-10 uppercase letters")
	}

	exists, err := s.repo.KeyExists(ctx, in.WorkspaceID, in.Key)
	if err != nil {
		return nil, fmt.Errorf("check key: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("project key %q already exists in workspace", in.Key)
	}

	p := &domain.Project{
		ID:          uuid.New(),
		WorkspaceID: in.WorkspaceID,
		Name:        in.Name,
		Key:         in.Key,
		CardCounter: 0,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	// Auto-create default board with columns
	_, err = s.boardSvc.Create(ctx, CreateBoardInput{
		ProjectID: p.ID,
		Name:      "Main Board",
	})
	if err != nil {
		return nil, fmt.Errorf("create default board: %w", err)
	}

	return p, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.Project, error) {
	return s.repo.ListByWorkspace(ctx, workspaceID)
}

type UpdateProjectInput struct {
	Name string
	Key  string
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, in UpdateProjectInput) (*domain.Project, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Key != "" && in.Key != p.Key {
		if !projectKeyRegex.MatchString(in.Key) {
			return nil, fmt.Errorf("project key must be 2-10 uppercase letters")
		}
		exists, err := s.repo.KeyExists(ctx, p.WorkspaceID, in.Key)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("project key %q already exists", in.Key)
		}
		p.Key = in.Key
	}
	if in.Name != "" {
		p.Name = in.Name
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
