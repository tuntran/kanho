package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/repository"
)

type WorkspaceService struct {
	repo *repository.WorkspaceRepo
}

func NewWorkspaceService(repo *repository.WorkspaceRepo) *WorkspaceService {
	return &WorkspaceService{repo: repo}
}

type CreateWorkspaceInput struct {
	Name        string
	Description string
	CreatedBy   uuid.UUID
}

func (s *WorkspaceService) Create(ctx context.Context, in CreateWorkspaceInput) (*domain.Workspace, error) {
	slug, err := domain.UniqueSlug(in.Name, func(slug string) (bool, error) {
		return s.repo.SlugExists(ctx, slug)
	})
	if err != nil {
		return nil, fmt.Errorf("generate slug: %w", err)
	}

	now := time.Now()
	w := &domain.Workspace{
		ID:          uuid.New(),
		Name:        in.Name,
		Slug:        slug,
		Description: in.Description,
		AccentColor: "#6366f1",
		CreatedBy:   in.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, w); err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}

	// Add creator as owner
	member := &domain.WorkspaceMember{
		WorkspaceID: w.ID,
		UserID:      in.CreatedBy,
		Role:        domain.RoleOwner,
		CreatedAt:   now,
	}
	if err := s.repo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("add owner member: %w", err)
	}

	return w, nil
}

func (s *WorkspaceService) GetBySlug(ctx context.Context, slug string) (*domain.Workspace, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *WorkspaceService) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Workspace, error) {
	return s.repo.ListByUser(ctx, userID)
}

type UpdateWorkspaceInput struct {
	Name        string
	Description string
	LogoURL     string
	AccentColor string
}

func (s *WorkspaceService) Update(ctx context.Context, id uuid.UUID, in UpdateWorkspaceInput) (*domain.Workspace, error) {
	w, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Name != "" && in.Name != w.Name {
		slug, err := domain.UniqueSlug(in.Name, func(slug string) (bool, error) {
			return s.repo.SlugExists(ctx, slug)
		})
		if err != nil {
			return nil, fmt.Errorf("generate slug: %w", err)
		}
		w.Name = in.Name
		w.Slug = slug
	}
	if in.Description != "" {
		w.Description = in.Description
	}
	if in.LogoURL != "" {
		w.LogoURL = in.LogoURL
	}
	if in.AccentColor != "" {
		w.AccentColor = in.AccentColor
	}

	if err := s.repo.Update(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WorkspaceService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Members

func (s *WorkspaceService) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]domain.WorkspaceMember, error) {
	return s.repo.ListMembers(ctx, workspaceID)
}

type InviteMemberInput struct {
	WorkspaceID uuid.UUID
	UserID      uuid.UUID
	Role        domain.Role
}

func (s *WorkspaceService) InviteMember(ctx context.Context, in InviteMemberInput) error {
	m := &domain.WorkspaceMember{
		WorkspaceID: in.WorkspaceID,
		UserID:      in.UserID,
		Role:        in.Role,
		CreatedAt:   time.Now(),
	}
	return s.repo.AddMember(ctx, m)
}

func (s *WorkspaceService) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role domain.Role) error {
	return s.repo.UpdateMemberRole(ctx, workspaceID, userID, role)
}

func (s *WorkspaceService) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	return s.repo.RemoveMember(ctx, workspaceID, userID)
}

func (s *WorkspaceService) GetMember(ctx context.Context, workspaceID, userID uuid.UUID) (*domain.WorkspaceMember, error) {
	return s.repo.GetMember(ctx, workspaceID, userID)
}
