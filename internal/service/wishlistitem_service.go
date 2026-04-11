package service

import (
	"context"
	"strings"
	"time"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/repository"
	"github.com/CodebyTecs/wishlist-service/pkg/uid"
)

const (
	minWishlistItemPriority = 1
	maxWishlistItemPriority = 5
)

type WishlistItemService interface {
	Create(ctx context.Context, userID, wishlistID string, input domain.CreateWishlistItemInput) (domain.WishlistItem, error)
	List(ctx context.Context, userID, wishlistID string) ([]domain.WishlistItem, error)
	GetByID(ctx context.Context, userID, wishlistID, itemID string) (domain.WishlistItem, error)
	UpdateByID(ctx context.Context, userID, wishlistID, itemID string, input domain.UpdateWishlistItemInput) (domain.WishlistItem, error)
	DeleteByID(ctx context.Context, userID, wishlistID, itemID string) error
}

type wishlistItemService struct {
	repo repository.WishlistItemRepository
}

func NewWishlistItemService(repo repository.WishlistItemRepository) WishlistItemService {
	return &wishlistItemService{repo: repo}
}

func (s *wishlistItemService) Create(ctx context.Context, userID, wishlistID string, input domain.CreateWishlistItemInput) (domain.WishlistItem, error) {
	if !isValidUserAndWishlistIDs(userID, wishlistID) {
		return domain.WishlistItem{}, domain.ErrInvalidRequest
	}
	if strings.TrimSpace(input.Name) == "" {
		return domain.WishlistItem{}, domain.ErrInvalidRequest
	}
	if input.Priority < minWishlistItemPriority || input.Priority > maxWishlistItemPriority {
		return domain.WishlistItem{}, domain.ErrInvalidRequest
	}

	itemID, err := uid.NewUUID()
	if err != nil {
		return domain.WishlistItem{}, err
	}

	item := domain.WishlistItem{
		ID:          itemID,
		WishlistID:  wishlistID,
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		ProductLink: strings.TrimSpace(input.ProductLink),
		Priority:    input.Priority,
		IsReserved:  false,
		CreatedAt:   time.Now().UTC(),
	}

	return s.repo.Create(ctx, userID, wishlistID, item)
}

func (s *wishlistItemService) List(ctx context.Context, userID, wishlistID string) ([]domain.WishlistItem, error) {
	if !isValidUserAndWishlistIDs(userID, wishlistID) {
		return nil, domain.ErrInvalidRequest
	}
	return s.repo.ListByWishlistAndUser(ctx, userID, wishlistID)
}

func (s *wishlistItemService) GetByID(ctx context.Context, userID, wishlistID, itemID string) (domain.WishlistItem, error) {
	if !isValidUserAndWishlistIDs(userID, wishlistID) || !uid.IsValidUUID(strings.TrimSpace(itemID)) {
		return domain.WishlistItem{}, domain.ErrInvalidRequest
	}
	return s.repo.GetByIDAndWishlistAndUser(ctx, strings.TrimSpace(userID), strings.TrimSpace(wishlistID), strings.TrimSpace(itemID))
}

func (s *wishlistItemService) UpdateByID(ctx context.Context, userID, wishlistID, itemID string, input domain.UpdateWishlistItemInput) (domain.WishlistItem, error) {
	if !isValidUserAndWishlistIDs(userID, wishlistID) || !uid.IsValidUUID(strings.TrimSpace(itemID)) {
		return domain.WishlistItem{}, domain.ErrInvalidRequest
	}
	if !input.UpdateName && !input.UpdateDescription && !input.UpdateProductLink && !input.UpdatePriority {
		return domain.WishlistItem{}, domain.ErrInvalidRequest
	}
	if input.UpdatePriority && (input.Priority < minWishlistItemPriority || input.Priority > maxWishlistItemPriority) {
		return domain.WishlistItem{}, domain.ErrInvalidRequest
	}

	update := domain.WishlistItemUpdate{
		Name:              strings.TrimSpace(input.Name),
		Description:       strings.TrimSpace(input.Description),
		ProductLink:       strings.TrimSpace(input.ProductLink),
		Priority:          input.Priority,
		UpdateName:        input.UpdateName,
		UpdateDescription: input.UpdateDescription,
		UpdateProductLink: input.UpdateProductLink,
		UpdatePriority:    input.UpdatePriority,
	}

	return s.repo.UpdateByIDAndWishlistAndUser(
		ctx,
		strings.TrimSpace(userID),
		strings.TrimSpace(wishlistID),
		strings.TrimSpace(itemID),
		update,
	)
}

func (s *wishlistItemService) DeleteByID(ctx context.Context, userID, wishlistID, itemID string) error {
	if !isValidUserAndWishlistIDs(userID, wishlistID) || !uid.IsValidUUID(strings.TrimSpace(itemID)) {
		return domain.ErrInvalidRequest
	}
	return s.repo.DeleteByIDAndWishlistAndUser(
		ctx,
		strings.TrimSpace(userID),
		strings.TrimSpace(wishlistID),
		strings.TrimSpace(itemID),
	)
}

func isValidUserAndWishlistIDs(userID, wishlistID string) bool {
	userID = strings.TrimSpace(userID)
	wishlistID = strings.TrimSpace(wishlistID)
	return uid.IsValidUUID(userID) && uid.IsValidUUID(wishlistID)
}
