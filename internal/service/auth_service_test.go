package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
)

type authUserRepoStub struct {
	usersByEmail  map[string]domain.User
	created       []domain.User
	createErr     error
	getByEmailErr error
	getByIDErr    error
}

func newAuthUserRepoStub() *authUserRepoStub {
	return &authUserRepoStub{
		usersByEmail: make(map[string]domain.User),
		created:      make([]domain.User, 0),
	}
}

func (s *authUserRepoStub) Create(_ context.Context, user domain.User) error {
	if s.createErr != nil {
		return s.createErr
	}
	if _, exists := s.usersByEmail[user.Email]; exists {
		return domain.ErrAlreadyExists
	}
	s.usersByEmail[user.Email] = user
	s.created = append(s.created, user)
	return nil
}

func (s *authUserRepoStub) GetByEmail(_ context.Context, email string) (domain.User, error) {
	if s.getByEmailErr != nil {
		return domain.User{}, s.getByEmailErr
	}
	user, ok := s.usersByEmail[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (s *authUserRepoStub) GetByID(_ context.Context, id string) (domain.User, error) {
	if s.getByIDErr != nil {
		return domain.User{}, s.getByIDErr
	}
	for _, user := range s.usersByEmail {
		if user.ID == id {
			return user, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

type tokenServiceStub struct {
	tokenToReturn string
	generateErr   error
}

func (s *tokenServiceStub) GenerateAccessToken(_ string) (string, error) {
	if s.generateErr != nil {
		return "", s.generateErr
	}
	return s.tokenToReturn, nil
}

func (s *tokenServiceStub) ParseAccessToken(_ string) (string, error) {
	return "", nil
}

func TestNewAuthService(t *testing.T) {
	repo := newAuthUserRepoStub()
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
}

func TestNormalizeEmail(t *testing.T) {
	got := normalizeEmail("  Test@Example.COM  ")
	if got != "test@example.com" {
		t.Fatalf("unexpected normalize result: %q", got)
	}
}

func TestAuthServiceRegisterInvalidRequest(t *testing.T) {
	repo := newAuthUserRepoStub()
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	_, err := svc.Register(context.Background(), "", "pass")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestAuthServiceRegisterAlreadyExists(t *testing.T) {
	repo := newAuthUserRepoStub()
	repo.usersByEmail["test@example.com"] = domain.User{ID: "1", Email: "test@example.com"}
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	_, err := svc.Register(context.Background(), "test@example.com", "qwerty123")
	if err != domain.ErrAlreadyExists {
		t.Fatalf("expected ErrAlreadyExists, got: %v", err)
	}
}

func TestAuthServiceRegisterGetByEmailUnexpectedError(t *testing.T) {
	errBoom := errors.New("boom")
	repo := newAuthUserRepoStub()
	repo.getByEmailErr = errBoom
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	_, err := svc.Register(context.Background(), "test@example.com", "qwerty123")
	if err != errBoom {
		t.Fatalf("expected boom error, got: %v", err)
	}
}

func TestAuthServiceRegisterBcryptError(t *testing.T) {
	repo := newAuthUserRepoStub()
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	longPassword := strings.Repeat("a", 73)
	_, err := svc.Register(context.Background(), "test@example.com", longPassword)
	if err == nil {
		t.Fatalf("expected bcrypt error for too long password")
	}
}

func TestAuthServiceRegisterCreateError(t *testing.T) {
	errBoom := errors.New("create error")

	repo := newAuthUserRepoStub()
	repo.createErr = errBoom
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	_, err := svc.Register(context.Background(), "test@example.com", "qwerty123")
	if err != errBoom {
		t.Fatalf("expected create error, got: %v", err)
	}
}

func TestAuthServiceRegisterTokenError(t *testing.T) {
	errBoom := errors.New("token error")

	repo := newAuthUserRepoStub()
	tokens := &tokenServiceStub{generateErr: errBoom}
	svc := NewAuthService(repo, tokens)

	_, err := svc.Register(context.Background(), "test@example.com", "qwerty123")
	if err != errBoom {
		t.Fatalf("expected token error, got: %v", err)
	}
}

func TestAuthServiceRegisterSuccess(t *testing.T) {
	repo := newAuthUserRepoStub()
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	token, err := svc.Register(context.Background(), "Test@Example.com", "qwerty123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "token-123" {
		t.Fatalf("unexpected token: %s", token)
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected one created user, got %d", len(repo.created))
	}

	createdUser := repo.created[0]
	if createdUser.PasswordHash == "qwerty123" {
		t.Fatalf("password hash must not equal plain password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(createdUser.PasswordHash), []byte("qwerty123")); err != nil {
		t.Fatalf("created hash does not match plain password: %v", err)
	}
}

func TestAuthServiceLoginInvalidRequest(t *testing.T) {
	repo := newAuthUserRepoStub()
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	_, err := svc.Login(context.Background(), "", "password")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestAuthServiceLoginUserNotFound(t *testing.T) {
	repo := newAuthUserRepoStub()
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	_, err := svc.Login(context.Background(), "user@example.com", "password")
	if err != domain.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestAuthServiceLoginGetByEmailUnexpectedError(t *testing.T) {
	errBoom := errors.New("repo error")
	repo := newAuthUserRepoStub()
	repo.getByEmailErr = errBoom
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	_, err := svc.Login(context.Background(), "user@example.com", "password")
	if err != errBoom {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

func TestAuthServiceLoginWrongPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("right-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	repo := newAuthUserRepoStub()
	repo.usersByEmail["user@example.com"] = domain.User{
		ID:           "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		Email:        "user@example.com",
		PasswordHash: string(passwordHash),
	}
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	_, err = svc.Login(context.Background(), "user@example.com", "wrong-password")
	if err != domain.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestAuthServiceLoginTokenError(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("right-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	errBoom := errors.New("token error")
	repo := newAuthUserRepoStub()
	repo.usersByEmail["user@example.com"] = domain.User{
		ID:           "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		Email:        "user@example.com",
		PasswordHash: string(passwordHash),
	}
	tokens := &tokenServiceStub{generateErr: errBoom}
	svc := NewAuthService(repo, tokens)

	_, err = svc.Login(context.Background(), "user@example.com", "right-password")
	if err != errBoom {
		t.Fatalf("expected token error, got: %v", err)
	}
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("right-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	repo := newAuthUserRepoStub()
	repo.usersByEmail["user@example.com"] = domain.User{
		ID:           "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		Email:        "user@example.com",
		PasswordHash: string(passwordHash),
	}
	tokens := &tokenServiceStub{tokenToReturn: "token-123"}
	svc := NewAuthService(repo, tokens)

	token, err := svc.Login(context.Background(), "user@example.com", "right-password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "token-123" {
		t.Fatalf("unexpected token: %s", token)
	}
}
