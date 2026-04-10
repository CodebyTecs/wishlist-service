package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/repository"
	"github.com/CodebyTecs/wishlist-service/pkg/uid"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	ParseAccessToken(token string) (string, error)
}

type authService struct {
	users  repository.UserRepository
	tokens TokenService
}

func NewAuthService(users repository.UserRepository, tokens TokenService) AuthService {
	return &authService{
		users:  users,
		tokens: tokens,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (string, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return "", domain.ErrInvalidRequest
	}

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return "", domain.ErrAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return "", err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userID, err := uid.NewUUID()
	if err != nil {
		return "", err
	}

	user := domain.User{
		ID:           userID,
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.users.Create(ctx, user); err != nil {
		return "", err
	}

	return s.tokens.GenerateAccessToken(userID)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return "", domain.ErrInvalidRequest
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrUnauthorized
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrUnauthorized
	}

	return s.tokens.GenerateAccessToken(user.ID)
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
