package service

import (
	"context"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/repository"
	"github.com/CodebyTecs/wishlist-service/pkg/uid"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (domain.User, error)
}

type userService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) UserService {
	return &userService{users: users}
}

func (s *userService) GetByID(ctx context.Context, id string) (domain.User, error) {
	if id == "" || !uid.IsValidUUID(id) {
		return domain.User{}, domain.ErrInvalidRequest
	}

	return s.users.GetByID(ctx, id)
}
