package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/auth"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/repository"
)

var (
	ErrEmailTaken        = errors.New("email already registered")
	ErrInvalidCredential = errors.New("invalid email or password")
	ErrTokenRevoked      = errors.New("refresh token revoked")
	ErrTokenExpired      = errors.New("refresh token expired")
	ErrTokenNotFound     = errors.New("refresh token not found")
	ErrUserNotFound      = errors.New("user not found")
)

type AuthService struct {
	users         *repository.UserRepository
	refreshTokens *repository.RefreshTokenRepository
	jwtSecret     string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewAuthService(
	users *repository.UserRepository,
	refreshTokens *repository.RefreshTokenRepository,
	jwtSecret string,
	accessTTL, refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		jwtSecret:     jwtSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (*domain.User, *auth.TokenPair, error) {
	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("checking email: %w", err)
	}
	if existing != nil {
		return nil, nil, ErrEmailTaken
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.users.Create(ctx, email, hash, name)
	if err != nil {
		return nil, nil, err
	}

	pair, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, *auth.TokenPair, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return nil, nil, ErrInvalidCredential
	}

	if err := auth.VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, nil, ErrInvalidCredential
	}

	pair, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshTokenStr string) (*domain.User, *auth.TokenPair, error) {
	claims, err := auth.ValidateToken(s.jwtSecret, refreshTokenStr)
	if err != nil {
		return nil, nil, ErrTokenExpired
	}

	jti := claims.ID
	if jti == "" {
		return nil, nil, ErrTokenNotFound
	}

	tokenHash := auth.HashJTI(jti)
	stored, err := s.refreshTokens.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, nil, fmt.Errorf("finding refresh token: %w", err)
	}
	if stored == nil {
		return nil, nil, ErrTokenNotFound
	}
	if stored.Revoked {
		// Possible token reuse attack — revoke all tokens for this user
		_ = s.refreshTokens.RevokeAllForUser(ctx, stored.UserID)
		return nil, nil, ErrTokenRevoked
	}

	// Revoke old token (rotation)
	if err := s.refreshTokens.Revoke(ctx, tokenHash); err != nil {
		return nil, nil, err
	}

	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	pair, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	claims, err := auth.ValidateToken(s.jwtSecret, refreshTokenStr)
	if err != nil {
		return nil // Already expired, nothing to revoke
	}

	if claims.ID == "" {
		return nil
	}

	tokenHash := auth.HashJTI(claims.ID)
	return s.refreshTokens.Revoke(ctx, tokenHash)
}

func (s *AuthService) RequiresSetup(ctx context.Context) (bool, error) {
	count, err := s.users.Count(ctx)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (s *AuthService) issueTokens(ctx context.Context, user *domain.User) (*auth.TokenPair, error) {
	pair, err := auth.GenerateTokenPair(s.jwtSecret, user.ID, user.Email, s.accessTTL, s.refreshTTL)
	if err != nil {
		return nil, err
	}

	tokenHash := auth.HashJTI(pair.RefreshJTI)
	expiresAt := time.Now().Add(s.refreshTTL)
	if err := s.refreshTokens.Store(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return nil, err
	}

	return pair, nil
}
