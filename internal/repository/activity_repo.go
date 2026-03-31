package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/domain"
)

type ActivityRepo struct {
	db *pgxpool.Pool
}

func NewActivityRepo(db *pgxpool.Pool) *ActivityRepo {
	return &ActivityRepo{db: db}
}

func (r *ActivityRepo) Create(ctx context.Context, a *domain.Activity) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO activities (id, workspace_id, actor_id, action, target_type, target_id, project_id, board_id, data, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		a.ID, a.WorkspaceID, a.ActorID, a.Action, a.TargetType, a.TargetID,
		a.ProjectID, a.BoardID, a.Data, a.CreatedAt,
	)
	return err
}

func (r *ActivityRepo) ListByCard(ctx context.Context, cardID uuid.UUID, limit, offset int) ([]domain.Activity, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(ctx,
		`SELECT a.id, a.workspace_id, a.actor_id, a.action, a.target_type, a.target_id,
		        COALESCE(a.project_id, '00000000-0000-0000-0000-000000000000'),
		        COALESCE(a.board_id, '00000000-0000-0000-0000-000000000000'),
		        a.data, a.created_at,
		        COALESCE(u.name,'')
		 FROM activities a
		 LEFT JOIN users u ON u.id = a.actor_id
		 WHERE a.target_id = $1
		 ORDER BY a.created_at DESC
		 LIMIT $2 OFFSET $3`, cardID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, scanActivity)
}

func scanActivity(row pgx.CollectableRow) (domain.Activity, error) {
	var a domain.Activity
	err := row.Scan(
		&a.ID, &a.WorkspaceID, &a.ActorID, &a.Action, &a.TargetType, &a.TargetID,
		&a.ProjectID, &a.BoardID, &a.Data, &a.CreatedAt,
		&a.ActorName,
	)
	return a, err
}
