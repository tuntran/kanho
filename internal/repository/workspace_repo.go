package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/domain"
)

type WorkspaceRepo struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepo(db *pgxpool.Pool) *WorkspaceRepo {
	return &WorkspaceRepo{db: db}
}

func (r *WorkspaceRepo) Create(ctx context.Context, w *domain.Workspace) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO workspaces (id, name, slug, description, logo_url, accent_color, created_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		w.ID, w.Name, w.Slug, w.Description, w.LogoURL, w.AccentColor, w.CreatedBy, w.CreatedAt, w.UpdatedAt,
	)
	return err
}

func (r *WorkspaceRepo) GetBySlug(ctx context.Context, slug string) (*domain.Workspace, error) {
	w := &domain.Workspace{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, slug, description, COALESCE(logo_url,''), COALESCE(accent_color,'#6366f1'), created_by, created_at, updated_at
		 FROM workspaces WHERE slug = $1`, slug,
	).Scan(&w.ID, &w.Name, &w.Slug, &w.Description, &w.LogoURL, &w.AccentColor, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *WorkspaceRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Workspace, error) {
	w := &domain.Workspace{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, slug, description, COALESCE(logo_url,''), COALESCE(accent_color,'#6366f1'), created_by, created_at, updated_at
		 FROM workspaces WHERE id = $1`, id,
	).Scan(&w.ID, &w.Name, &w.Slug, &w.Description, &w.LogoURL, &w.AccentColor, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *WorkspaceRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Workspace, error) {
	rows, err := r.db.Query(ctx,
		`SELECT w.id, w.name, w.slug, w.description, COALESCE(w.logo_url,''), COALESCE(w.accent_color,'#6366f1'), w.created_by, w.created_at, w.updated_at
		 FROM workspaces w
		 JOIN workspace_members wm ON wm.workspace_id = w.id
		 WHERE wm.user_id = $1
		 ORDER BY w.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Workspace, error) {
		var w domain.Workspace
		err := row.Scan(&w.ID, &w.Name, &w.Slug, &w.Description, &w.LogoURL, &w.AccentColor, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
		return w, err
	})
}

func (r *WorkspaceRepo) Update(ctx context.Context, w *domain.Workspace) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE workspaces SET name=$1, slug=$2, description=$3, logo_url=$4, accent_color=$5, updated_at=now()
		 WHERE id=$6`,
		w.Name, w.Slug, w.Description, w.LogoURL, w.AccentColor, w.ID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("workspace not found")
	}
	return nil
}

func (r *WorkspaceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM workspaces WHERE id = $1`, id)
	return err
}

func (r *WorkspaceRepo) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM workspaces WHERE slug = $1)`, slug).Scan(&exists)
	return exists, err
}

// Members

func (r *WorkspaceRepo) AddMember(ctx context.Context, m *domain.WorkspaceMember) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO workspace_members (workspace_id, user_id, role, created_at)
		 VALUES ($1, $2, $3, $4)`,
		m.WorkspaceID, m.UserID, m.Role, m.CreatedAt,
	)
	return err
}

func (r *WorkspaceRepo) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]domain.WorkspaceMember, error) {
	rows, err := r.db.Query(ctx,
		`SELECT wm.workspace_id, wm.user_id, wm.role, wm.created_at, u.name, u.email
		 FROM workspace_members wm
		 JOIN users u ON u.id = wm.user_id
		 WHERE wm.workspace_id = $1
		 ORDER BY wm.created_at`, workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.WorkspaceMember, error) {
		var m domain.WorkspaceMember
		err := row.Scan(&m.WorkspaceID, &m.UserID, &m.Role, &m.CreatedAt, &m.UserName, &m.UserEmail)
		return m, err
	})
}

func (r *WorkspaceRepo) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role domain.Role) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE workspace_members SET role = $1 WHERE workspace_id = $2 AND user_id = $3`,
		role, workspaceID, userID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("member not found")
	}
	return nil
}

func (r *WorkspaceRepo) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID, userID,
	)
	return err
}

func (r *WorkspaceRepo) GetMember(ctx context.Context, workspaceID, userID uuid.UUID) (*domain.WorkspaceMember, error) {
	m := &domain.WorkspaceMember{}
	err := r.db.QueryRow(ctx,
		`SELECT workspace_id, user_id, role, created_at
		 FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID, userID,
	).Scan(&m.WorkspaceID, &m.UserID, &m.Role, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}
