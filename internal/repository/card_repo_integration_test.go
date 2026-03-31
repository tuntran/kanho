package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/repository"
	"github.com/tungtran/kanho/internal/repository/testutil"
)

func TestCardRepo_CreateAndRead(t *testing.T) {
	pool := testutil.NewTestPool(t)

	cardRepo := repository.NewCardRepo(pool)

	// Create prerequisite data
	ctx := context.Background()
	userID := uuid.New()
	wsID := uuid.New()
	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	now := time.Now()

	// Insert user
	_, err := pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5)`,
		userID, "integration@test.com", "hash", "Test User", now)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	// Insert workspace
	_, err = pool.Exec(ctx,
		`INSERT INTO workspaces (id, name, slug, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5)`,
		wsID, "Test WS", "test-ws", userID, now)
	if err != nil {
		t.Fatalf("insert workspace: %v", err)
	}

	// Insert project
	_, err = pool.Exec(ctx,
		`INSERT INTO projects (id, workspace_id, name, key, slug, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $7)`,
		projectID, wsID, "Test Project", "TST", "test-project", userID, now)
	if err != nil {
		t.Fatalf("insert project: %v", err)
	}

	// Insert board
	_, err = pool.Exec(ctx,
		`INSERT INTO boards (id, project_id, name, slug, visibility, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $6)`,
		boardID, projectID, "Test Board", "test-board", "private", now)
	if err != nil {
		t.Fatalf("insert board: %v", err)
	}

	// Insert column
	_, err = pool.Exec(ctx,
		`INSERT INTO columns (id, board_id, name, position, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5)`,
		columnID, boardID, "Todo", "a", now)
	if err != nil {
		t.Fatalf("insert column: %v", err)
	}

	// Create card
	card := &domain.Card{
		ID:        uuid.New(),
		ProjectID: projectID,
		BoardID:   boardID,
		ColumnID:  columnID,
		Title:     "Integration Test Card",
		Priority:  domain.PriorityMedium,
		Position:  "m",
		CreatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := cardRepo.Create(ctx, card); err != nil {
		t.Fatalf("Create card: %v", err)
	}

	if card.CardNumber == 0 {
		t.Error("CardNumber should be assigned by Create")
	}

	// Read back
	got, err := cardRepo.GetByID(ctx, card.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}

	if got.Title != "Integration Test Card" {
		t.Errorf("expected title %q, got %q", "Integration Test Card", got.Title)
	}
	if got.CardNumber != card.CardNumber {
		t.Errorf("expected card_number %d, got %d", card.CardNumber, got.CardNumber)
	}
}
