package service

import (
	"context"
	"strings"
	"time"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/repository"
	"github.com/CodebyTecs/wishlist-service/pkg/uid"
)

type WishlistService interface {
	Create(ctx context.Context, userID string, input domain.CreateWishlistInput) (domain.Wishlist, error)
	ListByUserID(ctx context.Context, userID string) ([]domain.Wishlist, error)
	GetByID(ctx context.Context, userID, wishlistID string) (domain.Wishlist, error)
	UpdateByID(ctx context.Context, userID, wishlistID string, input domain.UpdateWishlistInput) (domain.Wishlist, error)
	DeleteByID(ctx context.Context, userID, wishlistID string) error
}

type wishlistService struct {
	repo repository.WishlistRepository
}

func NewWishlistService(repo repository.WishlistRepository) WishlistService {
	return &wishlistService{repo: repo}
}

func (s *wishlistService) Create(ctx context.Context, userID string, input domain.CreateWishlistInput) (domain.Wishlist, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" || !uid.IsValidUUID(userID) {
		return domain.Wishlist{}, domain.ErrInvalidRequest
	}
	if strings.TrimSpace(input.Name) == "" || input.EventDate.IsZero() {
		return domain.Wishlist{}, domain.ErrInvalidRequest
	}

	wishlistID, err := uid.NewUUID()
	if err != nil {
		return domain.Wishlist{}, err
	}
	publicToken, err := uid.NewUUID()
	if err != nil {
		return domain.Wishlist{}, err
	}

	wishlist := domain.Wishlist{
		ID:          wishlistID,
		UserID:      userID,
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		EventDate:   input.EventDate,
		PublicToken: publicToken,
		CreatedAt:   time.Now().UTC(),
	}

	return s.repo.Create(ctx, wishlist)
}

func (s *wishlistService) ListByUserID(ctx context.Context, userID string) ([]domain.Wishlist, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" || !uid.IsValidUUID(userID) {
		return nil, domain.ErrInvalidRequest
	}

	return s.repo.ListByUserID(ctx, userID)
}

func (s *wishlistService) GetByID(ctx context.Context, userID, wishlistID string) (domain.Wishlist, error) {
	userID = strings.TrimSpace(userID)
	wishlistID = strings.TrimSpace(wishlistID)
	if userID == "" || wishlistID == "" || !uid.IsValidUUID(userID) || !uid.IsValidUUID(wishlistID) {
		return domain.Wishlist{}, domain.ErrInvalidRequest
	}

	return s.repo.GetByIDAndUserID(ctx, wishlistID, userID)
}

func (s *wishlistService) UpdateByID(ctx context.Context, userID, wishlistID string, input domain.UpdateWishlistInput) (domain.Wishlist, error) {
	userID = strings.TrimSpace(userID)
	wishlistID = strings.TrimSpace(wishlistID)
	if userID == "" || wishlistID == "" || !uid.IsValidUUID(userID) || !uid.IsValidUUID(wishlistID) {
		return domain.Wishlist{}, domain.ErrInvalidRequest
	}
	if !input.UpdateName && !input.UpdateDescription && !input.UpdateEventDate {
		return domain.Wishlist{}, domain.ErrInvalidRequest
	}

	return s.repo.UpdateByIDAndUserID(ctx, wishlistID, userID, domain.WishlistUpdate{
		Name:              strings.TrimSpace(input.Name),
		Description:       strings.TrimSpace(input.Description),
		EventDate:         input.EventDate,
		UpdateName:        input.UpdateName,
		UpdateDescription: input.UpdateDescription,
		UpdateEventDate:   input.UpdateEventDate,
	})
}

func (s *wishlistService) DeleteByID(ctx context.Context, userID, wishlistID string) error {
	userID = strings.TrimSpace(userID)
	wishlistID = strings.TrimSpace(wishlistID)
	if userID == "" || wishlistID == "" || !uid.IsValidUUID(userID) || !uid.IsValidUUID(wishlistID) {
		return domain.ErrInvalidRequest
	}

	return s.repo.DeleteByIDAndUserID(ctx, wishlistID, userID)
}
