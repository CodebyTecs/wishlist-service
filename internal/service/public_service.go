package service

import (
	"context"
	"strings"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/pkg/uid"
)

type publicWishlistRepository interface {
	GetByPublicToken(ctx context.Context, publicToken string) (domain.Wishlist, error)
}

type publicWishlistItemRepository interface {
	ListByPublicToken(ctx context.Context, publicToken string) ([]domain.WishlistItem, error)
	ReserveByPublicTokenAndItemID(ctx context.Context, publicToken, itemID string) error
}

type PublicService interface {
	GetWishlistByToken(ctx context.Context, publicToken string) (domain.Wishlist, []domain.WishlistItem, error)
	ReserveItem(ctx context.Context, publicToken, itemID string) error
}

type publicService struct {
	wishlists publicWishlistRepository
	items     publicWishlistItemRepository
}

func NewPublicService(wishlists publicWishlistRepository, items publicWishlistItemRepository) PublicService {
	return &publicService{
		wishlists: wishlists,
		items:     items,
	}
}

func (s *publicService) GetWishlistByToken(ctx context.Context, publicToken string) (domain.Wishlist, []domain.WishlistItem, error) {
	publicToken = strings.TrimSpace(publicToken)
	if publicToken == "" {
		return domain.Wishlist{}, nil, domain.ErrInvalidRequest
	}

	wishlist, err := s.wishlists.GetByPublicToken(ctx, publicToken)
	if err != nil {
		return domain.Wishlist{}, nil, err
	}

	items, err := s.items.ListByPublicToken(ctx, publicToken)
	if err != nil {
		return domain.Wishlist{}, nil, err
	}

	return wishlist, items, nil
}

func (s *publicService) ReserveItem(ctx context.Context, publicToken, itemID string) error {
	publicToken = strings.TrimSpace(publicToken)
	itemID = strings.TrimSpace(itemID)
	if publicToken == "" || !uid.IsValidUUID(itemID) {
		return domain.ErrInvalidRequest
	}

	return s.items.ReserveByPublicTokenAndItemID(ctx, publicToken, itemID)
}
