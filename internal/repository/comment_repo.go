package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/domain"
)

type CommentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, c *domain.Comment) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO comments (id, card_id, author_id, content, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		c.ID, c.CardID, c.AuthorID, c.Content, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CommentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	c := &domain.Comment{}
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.card_id, c.author_id, c.content, c.created_at, c.updated_at,
		        COALESCE(u.name,'')
		 FROM comments c
		 LEFT JOIN users u ON u.id = c.author_id
		 WHERE c.id = $1`, id,
	).Scan(&c.ID, &c.CardID, &c.AuthorID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.AuthorName)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CommentRepo) ListByCard(ctx context.Context, cardID uuid.UUID, limit, offset int) ([]domain.Comment, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.card_id, c.author_id, c.content, c.created_at, c.updated_at,
		        COALESCE(u.name,'')
		 FROM comments c
		 LEFT JOIN users u ON u.id = c.author_id
		 WHERE c.card_id = $1
		 ORDER BY c.created_at DESC
		 LIMIT $2 OFFSET $3`, cardID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, scanComment)
}

func (r *CommentRepo) Update(ctx context.Context, id uuid.UUID, content string) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE comments SET content=$1, updated_at=now() WHERE id=$2`,
		content, id,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("comment not found")
	}
	return nil
}

func (r *CommentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM comments WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("comment not found")
	}
	return nil
}

func scanComment(row pgx.CollectableRow) (domain.Comment, error) {
	var c domain.Comment
	err := row.Scan(&c.ID, &c.CardID, &c.AuthorID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.AuthorName)
	return c, err
}
