package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/domain"
)

type BoardRepo struct {
	db *pgxpool.Pool
}

func NewBoardRepo(db *pgxpool.Pool) *BoardRepo {
	return &BoardRepo{db: db}
}

func (r *BoardRepo) Create(ctx context.Context, b *domain.Board) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO boards (id, project_id, name, slug, description, visibility, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		b.ID, b.ProjectID, b.Name, b.Slug, b.Description, b.Visibility, b.CreatedAt, b.UpdatedAt,
	)
	return err
}

func (r *BoardRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Board, error) {
	b := &domain.Board{}
	err := r.db.QueryRow(ctx,
		`SELECT id, project_id, name, slug, COALESCE(description,''), visibility, created_at, updated_at
		 FROM boards WHERE id = $1`, id,
	).Scan(&b.ID, &b.ProjectID, &b.Name, &b.Slug, &b.Description, &b.Visibility, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BoardRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.Board, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, project_id, name, slug, COALESCE(description,''), visibility, created_at, updated_at
		 FROM boards WHERE project_id = $1
		 ORDER BY created_at`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Board, error) {
		var b domain.Board
		err := row.Scan(&b.ID, &b.ProjectID, &b.Name, &b.Slug, &b.Description, &b.Visibility, &b.CreatedAt, &b.UpdatedAt)
		return b, err
	})
}

func (r *BoardRepo) Update(ctx context.Context, b *domain.Board) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE boards SET name=$1, slug=$2, description=$3, visibility=$4, updated_at=now()
		 WHERE id=$5`,
		b.Name, b.Slug, b.Description, b.Visibility, b.ID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("board not found")
	}
	return nil
}

func (r *BoardRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM boards WHERE id = $1`, id)
	return err
}

func (r *BoardRepo) SlugExists(ctx context.Context, projectID uuid.UUID, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM boards WHERE project_id = $1 AND slug = $2)`,
		projectID, slug,
	).Scan(&exists)
	return exists, err
}
