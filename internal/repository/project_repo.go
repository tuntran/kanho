package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/domain"
)

type ProjectRepo struct {
	db *pgxpool.Pool
}

func NewProjectRepo(db *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(ctx context.Context, p *domain.Project) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO projects (id, workspace_id, name, key, card_counter, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		p.ID, p.WorkspaceID, p.Name, p.Key, p.CardCounter, p.CreatedAt,
	)
	return err
}

func (r *ProjectRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	p := &domain.Project{}
	err := r.db.QueryRow(ctx,
		`SELECT id, workspace_id, name, key, card_counter, created_at
		 FROM projects WHERE id = $1`, id,
	).Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Key, &p.CardCounter, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProjectRepo) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.Project, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, workspace_id, name, key, card_counter, created_at
		 FROM projects WHERE workspace_id = $1
		 ORDER BY created_at DESC`, workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Project, error) {
		var p domain.Project
		err := row.Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Key, &p.CardCounter, &p.CreatedAt)
		return p, err
	})
}

func (r *ProjectRepo) Update(ctx context.Context, p *domain.Project) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE projects SET name = $1, key = $2 WHERE id = $3`,
		p.Name, p.Key, p.ID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}

func (r *ProjectRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	return err
}

func (r *ProjectRepo) KeyExists(ctx context.Context, workspaceID uuid.UUID, key string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM projects WHERE workspace_id = $1 AND key = $2)`,
		workspaceID, key,
	).Scan(&exists)
	return exists, err
}
