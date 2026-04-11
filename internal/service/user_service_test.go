package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
)

type userRepoStub struct {
	getByIDUser domain.User
	getByIDErr  error
}

func (s *userRepoStub) Create(_ context.Context, _ domain.User) error {
	return nil
}

func (s *userRepoStub) GetByEmail(_ context.Context, _ string) (domain.User, error) {
	return domain.User{}, domain.ErrNotFound
}

func (s *userRepoStub) GetByID(_ context.Context, _ string) (domain.User, error) {
	if s.getByIDErr != nil {
		return domain.User{}, s.getByIDErr
	}
	return s.getByIDUser, nil
}

func TestNewUserService(t *testing.T) {
	svc := NewUserService(&userRepoStub{})
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
}

func TestUserServiceGetByIDInvalidRequest(t *testing.T) {
	svc := NewUserService(&userRepoStub{})

	_, err := svc.GetByID(context.Background(), "bad-id")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestUserServiceGetByIDPassThroughError(t *testing.T) {
	errBoom := errors.New("repo error")
	svc := NewUserService(&userRepoStub{getByIDErr: errBoom})

	_, err := svc.GetByID(context.Background(), "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7")
	if err != errBoom {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

func TestUserServiceGetByIDSuccess(t *testing.T) {
	now := time.Now().UTC()
	svc := NewUserService(&userRepoStub{
		getByIDUser: domain.User{
			ID:        "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
			Email:     "user@example.com",
			CreatedAt: now,
		},
	})

	user, err := svc.GetByID(context.Background(), "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("unexpected user: %+v", user)
	}
}
