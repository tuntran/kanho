package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/domain"
)

type CardRepo struct {
	db *pgxpool.Pool
}

func NewCardRepo(db *pgxpool.Pool) *CardRepo {
	return &CardRepo{db: db}
}

// Create inserts a card and assigns a sequential card_number via next_card_number().
func (r *CardRepo) Create(ctx context.Context, c *domain.Card) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get next card number atomically
	var number int
	err = tx.QueryRow(ctx, `SELECT next_card_number($1)`, c.ProjectID).Scan(&number)
	if err != nil {
		return fmt.Errorf("next_card_number: %w", err)
	}
	c.CardNumber = number

	_, err = tx.Exec(ctx,
		`INSERT INTO cards (id, project_id, board_id, column_id, card_number, title, description, priority, position, start_date, due_date, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		c.ID, c.ProjectID, c.BoardID, c.ColumnID, c.CardNumber,
		c.Title, c.Description, c.Priority, c.Position,
		c.StartDate, c.DueDate, c.CreatedBy, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert card: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *CardRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Card, error) {
	c := &domain.Card{}
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.project_id, c.board_id, c.column_id, c.card_number,
		        c.title, COALESCE(c.description,''), c.priority, c.position,
		        c.start_date, c.due_date, c.archived_at,
		        c.created_by, c.created_at, c.updated_at,
		        p.key
		 FROM cards c
		 JOIN projects p ON p.id = c.project_id
		 WHERE c.id = $1`, id,
	).Scan(
		&c.ID, &c.ProjectID, &c.BoardID, &c.ColumnID, &c.CardNumber,
		&c.Title, &c.Description, &c.Priority, &c.Position,
		&c.StartDate, &c.DueDate, &c.ArchivedAt,
		&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
		&c.ReadableID,
	)
	if err != nil {
		return nil, err
	}
	c.ReadableID = c.ReadableID + "-" + fmt.Sprint(c.CardNumber)
	return c, nil
}

func (r *CardRepo) ListByColumn(ctx context.Context, columnID uuid.UUID) ([]domain.Card, error) {
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.project_id, c.board_id, c.column_id, c.card_number,
		        c.title, COALESCE(c.description,''), c.priority, c.position,
		        c.start_date, c.due_date, c.archived_at,
		        c.created_by, c.created_at, c.updated_at,
		        p.key
		 FROM cards c
		 JOIN projects p ON p.id = c.project_id
		 WHERE c.column_id = $1 AND c.archived_at IS NULL
		 ORDER BY c.position`, columnID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, scanCard)
}

func (r *CardRepo) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]domain.Card, error) {
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.project_id, c.board_id, c.column_id, c.card_number,
		        c.title, COALESCE(c.description,''), c.priority, c.position,
		        c.start_date, c.due_date, c.archived_at,
		        c.created_by, c.created_at, c.updated_at,
		        p.key
		 FROM cards c
		 JOIN projects p ON p.id = c.project_id
		 WHERE c.board_id = $1 AND c.archived_at IS NULL
		 ORDER BY c.position`, boardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, scanCard)
}

func (r *CardRepo) Update(ctx context.Context, c *domain.Card) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE cards SET title=$1, description=$2, priority=$3,
		        start_date=$4, due_date=$5, updated_at=now()
		 WHERE id=$6`,
		c.Title, c.Description, c.Priority,
		c.StartDate, c.DueDate, c.ID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("card not found")
	}
	return nil
}

// Move updates column_id and position in a single statement.
func (r *CardRepo) Move(ctx context.Context, cardID, columnID uuid.UUID, position string) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE cards SET column_id=$1, position=$2, updated_at=now() WHERE id=$3`,
		columnID, position, cardID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("card not found")
	}
	return nil
}

func (r *CardRepo) Archive(ctx context.Context, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE cards SET archived_at=now(), updated_at=now() WHERE id=$1`, id,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("card not found")
	}
	return nil
}

func (r *CardRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cards WHERE id=$1`, id)
	return err
}

// Positions returns sorted positions of non-archived cards in a column.
func (r *CardRepo) Positions(ctx context.Context, columnID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(ctx,
		`SELECT position FROM cards
		 WHERE column_id=$1 AND archived_at IS NULL
		 ORDER BY position`, columnID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var positions []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, rows.Err()
}

// AddAssignee adds a user as card assignee.
func (r *CardRepo) AddAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO card_assignees (card_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		cardID, userID,
	)
	return err
}

func (r *CardRepo) RemoveAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM card_assignees WHERE card_id=$1 AND user_id=$2`,
		cardID, userID,
	)
	return err
}

func (r *CardRepo) AddLabel(ctx context.Context, cardID, labelID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO card_labels (card_id, label_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		cardID, labelID,
	)
	return err
}

func (r *CardRepo) RemoveLabel(ctx context.Context, cardID, labelID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM card_labels WHERE card_id=$1 AND label_id=$2`,
		cardID, labelID,
	)
	return err
}

// SearchResult holds a card search result with board/column context.
type SearchResult struct {
	CardID     uuid.UUID `json:"card_id"`
	CardNumber int       `json:"card_number"`
	ReadableID string    `json:"readable_id"`
	Title      string    `json:"title"`
	BoardID    uuid.UUID `json:"board_id"`
	BoardName  string    `json:"board_name"`
	ColumnName string    `json:"column_name"`
}

// SearchFullText performs PostgreSQL full-text search on cards within a workspace.
func (r *CardRepo) SearchFullText(ctx context.Context, workspaceID uuid.UUID, query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.card_number, p.key, c.title, c.board_id,
		        b.name AS board_name, col.name AS column_name
		 FROM cards c
		 JOIN projects p ON p.id = c.project_id
		 JOIN boards b ON b.id = c.board_id
		 JOIN columns col ON col.id = c.column_id
		 WHERE p.workspace_id = $1
		   AND c.archived_at IS NULL
		   AND c.search_vector @@ plainto_tsquery('english', $2)
		 ORDER BY ts_rank(c.search_vector, plainto_tsquery('english', $2)) DESC
		 LIMIT $3`,
		workspaceID, query, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var sr SearchResult
		var projectKey string
		if err := rows.Scan(&sr.CardID, &sr.CardNumber, &projectKey, &sr.Title, &sr.BoardID, &sr.BoardName, &sr.ColumnName); err != nil {
			return nil, err
		}
		sr.ReadableID = projectKey + "-" + fmt.Sprint(sr.CardNumber)
		results = append(results, sr)
	}
	return results, rows.Err()
}

func scanCard(row pgx.CollectableRow) (domain.Card, error) {
	var c domain.Card
	var projectKey string
	err := row.Scan(
		&c.ID, &c.ProjectID, &c.BoardID, &c.ColumnID, &c.CardNumber,
		&c.Title, &c.Description, &c.Priority, &c.Position,
		&c.StartDate, &c.DueDate, &c.ArchivedAt,
		&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
		&projectKey,
	)
	c.ReadableID = projectKey + "-" + fmt.Sprint(c.CardNumber)
	return c, err
}
