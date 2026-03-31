package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/domain"
)

type ColumnRepo struct {
	db *pgxpool.Pool
}

func NewColumnRepo(db *pgxpool.Pool) *ColumnRepo {
	return &ColumnRepo{db: db}
}

func (r *ColumnRepo) Create(ctx context.Context, c *domain.Column) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO columns (id, board_id, name, position, color, wip_limit, is_done, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.BoardID, c.Name, c.Position, c.Color, c.WIPLimit, c.IsDone, c.CreatedAt,
	)
	return err
}

func (r *ColumnRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Column, error) {
	c := &domain.Column{}
	err := r.db.QueryRow(ctx,
		`SELECT id, board_id, name, position, COALESCE(color,''), wip_limit, is_done, created_at
		 FROM columns WHERE id = $1`, id,
	).Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &c.Color, &c.WIPLimit, &c.IsDone, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ColumnRepo) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]domain.Column, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, board_id, name, position, COALESCE(color,''), wip_limit, is_done, created_at
		 FROM columns WHERE board_id = $1
		 ORDER BY position`, boardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Column, error) {
		var c domain.Column
		err := row.Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &c.Color, &c.WIPLimit, &c.IsDone, &c.CreatedAt)
		return c, err
	})
}

func (r *ColumnRepo) Update(ctx context.Context, c *domain.Column) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE columns SET name=$1, position=$2, color=$3, wip_limit=$4, is_done=$5
		 WHERE id=$6`,
		c.Name, c.Position, c.Color, c.WIPLimit, c.IsDone, c.ID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("column not found")
	}
	return nil
}

func (r *ColumnRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM columns WHERE id = $1`, id)
	return err
}

// BatchUpdatePositions updates positions for multiple columns in a transaction.
func (r *ColumnRepo) BatchUpdatePositions(ctx context.Context, updates []domain.Column) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, u := range updates {
		_, err := tx.Exec(ctx,
			`UPDATE columns SET position = $1 WHERE id = $2`,
			u.Position, u.ID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
