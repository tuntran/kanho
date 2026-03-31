package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepo
}

func NewCommentService(commentRepo *repository.CommentRepo) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

type CreateCommentInput struct {
	CardID   uuid.UUID
	AuthorID uuid.UUID
	Content  string
}

func (s *CommentService) Create(ctx context.Context, in CreateCommentInput) (*domain.Comment, error) {
	if in.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	now := time.Now()
	c := &domain.Comment{
		ID:        uuid.New(),
		CardID:    in.CardID,
		AuthorID:  in.AuthorID,
		Content:   in.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.commentRepo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	// Re-fetch to get author name
	return s.commentRepo.GetByID(ctx, c.ID)
}

func (s *CommentService) ListByCard(ctx context.Context, cardID uuid.UUID, limit, offset int) ([]domain.Comment, error) {
	return s.commentRepo.ListByCard(ctx, cardID, limit, offset)
}

func (s *CommentService) Update(ctx context.Context, commentID, actorID uuid.UUID, content string) (*domain.Comment, error) {
	c, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}

	if c.AuthorID != actorID {
		return nil, fmt.Errorf("only the author can edit this comment")
	}

	if err := s.commentRepo.Update(ctx, commentID, content); err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}

	return s.commentRepo.GetByID(ctx, commentID)
}

func (s *CommentService) Delete(ctx context.Context, commentID, actorID uuid.UUID) error {
	c, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	// Only author can delete their own comment
	if c.AuthorID != actorID {
		return fmt.Errorf("only the author can delete this comment")
	}

	return s.commentRepo.Delete(ctx, commentID)
}
