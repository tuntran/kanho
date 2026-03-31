package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/domain"
)

type LabelRepo struct {
	db *pgxpool.Pool
}

func NewLabelRepo(db *pgxpool.Pool) *LabelRepo {
	return &LabelRepo{db: db}
}

func (r *LabelRepo) Create(ctx context.Context, l *domain.Label) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO labels (id, workspace_id, name, color, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		l.ID, l.WorkspaceID, l.Name, l.Color, l.CreatedAt,
	)
	return err
}

func (r *LabelRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Label, error) {
	l := &domain.Label{}
	err := r.db.QueryRow(ctx,
		`SELECT id, workspace_id, name, color, created_at
		 FROM labels WHERE id = $1`, id,
	).Scan(&l.ID, &l.WorkspaceID, &l.Name, &l.Color, &l.CreatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *LabelRepo) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.Label, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, workspace_id, name, color, created_at
		 FROM labels WHERE workspace_id = $1
		 ORDER BY name`, workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Label, error) {
		var l domain.Label
		err := row.Scan(&l.ID, &l.WorkspaceID, &l.Name, &l.Color, &l.CreatedAt)
		return l, err
	})
}

func (r *LabelRepo) Update(ctx context.Context, l *domain.Label) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE labels SET name = $1, color = $2 WHERE id = $3`,
		l.Name, l.Color, l.ID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("label not found")
	}
	return nil
}

func (r *LabelRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM labels WHERE id = $1`, id)
	return err
}
