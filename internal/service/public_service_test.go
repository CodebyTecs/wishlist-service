package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
)

type publicWishlistRepoStub struct {
	wishlist domain.Wishlist
	err      error
}

func (s *publicWishlistRepoStub) GetByPublicToken(_ context.Context, _ string) (domain.Wishlist, error) {
	return s.wishlist, s.err
}

type publicWishlistItemRepoStub struct {
	items      []domain.WishlistItem
	listErr    error
	reserveErr error
}

func (s *publicWishlistItemRepoStub) ListByPublicToken(_ context.Context, _ string) ([]domain.WishlistItem, error) {
	return s.items, s.listErr
}

func (s *publicWishlistItemRepoStub) ReserveByPublicTokenAndItemID(_ context.Context, _, _ string) error {
	return s.reserveErr
}

func TestNewPublicService(t *testing.T) {
	svc := NewPublicService(&publicWishlistRepoStub{}, &publicWishlistItemRepoStub{})
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
}

func TestPublicServiceGetWishlistByTokenInvalidRequest(t *testing.T) {
	svc := NewPublicService(&publicWishlistRepoStub{}, &publicWishlistItemRepoStub{})

	_, _, err := svc.GetWishlistByToken(context.Background(), "   ")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestPublicServiceGetWishlistByTokenWishlistError(t *testing.T) {
	errBoom := errors.New("wishlist error")
	svc := NewPublicService(
		&publicWishlistRepoStub{err: errBoom},
		&publicWishlistItemRepoStub{},
	)

	_, _, err := svc.GetWishlistByToken(context.Background(), "token-1")
	if err != errBoom {
		t.Fatalf("expected wishlist error, got: %v", err)
	}
}

func TestPublicServiceGetWishlistByTokenItemsError(t *testing.T) {
	errBoom := errors.New("items error")
	now := time.Now().UTC()
	svc := NewPublicService(
		&publicWishlistRepoStub{
			wishlist: domain.Wishlist{
				ID:          "6ef4f0a8-88ea-4ee8-911c-b5487f312380",
				Name:        "Birthday",
				Description: "My gifts",
				EventDate:   now,
				PublicToken: "token-1",
			},
		},
		&publicWishlistItemRepoStub{listErr: errBoom},
	)

	_, _, err := svc.GetWishlistByToken(context.Background(), "token-1")
	if err != errBoom {
		t.Fatalf("expected items error, got: %v", err)
	}
}

func TestPublicServiceReserveItemAlreadyReserved(t *testing.T) {
	svc := NewPublicService(
		&publicWishlistRepoStub{},
		&publicWishlistItemRepoStub{reserveErr: domain.ErrAlreadyReserved},
	)

	err := svc.ReserveItem(context.Background(), "public-token", "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7")
	if err != domain.ErrAlreadyReserved {
		t.Fatalf("expected ErrAlreadyReserved, got: %v", err)
	}
}

func TestPublicServiceReserveItemInvalidRequest(t *testing.T) {
	svc := NewPublicService(
		&publicWishlistRepoStub{},
		&publicWishlistItemRepoStub{},
	)

	err := svc.ReserveItem(context.Background(), "public-token", "not-a-uuid")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestPublicServiceReserveItemInvalidRequestEmptyToken(t *testing.T) {
	svc := NewPublicService(
		&publicWishlistRepoStub{},
		&publicWishlistItemRepoStub{},
	)

	err := svc.ReserveItem(context.Background(), "   ", "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestPublicServiceReserveItemRepoError(t *testing.T) {
	errBoom := errors.New("reserve error")
	svc := NewPublicService(
		&publicWishlistRepoStub{},
		&publicWishlistItemRepoStub{reserveErr: errBoom},
	)

	err := svc.ReserveItem(context.Background(), "public-token", "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7")
	if err != errBoom {
		t.Fatalf("expected reserve error, got: %v", err)
	}
}

func TestPublicServiceReserveItemSuccess(t *testing.T) {
	svc := NewPublicService(
		&publicWishlistRepoStub{},
		&publicWishlistItemRepoStub{},
	)

	err := svc.ReserveItem(context.Background(), "public-token", "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPublicServiceGetWishlistByTokenSuccess(t *testing.T) {
	now := time.Now().UTC()
	svc := NewPublicService(
		&publicWishlistRepoStub{
			wishlist: domain.Wishlist{
				ID:          "6ef4f0a8-88ea-4ee8-911c-b5487f312380",
				Name:        "Birthday",
				Description: "My gifts",
				EventDate:   now,
				PublicToken: "token-1",
			},
		},
		&publicWishlistItemRepoStub{
			items: []domain.WishlistItem{
				{
					ID:         "19f3638e-3274-4c99-8a34-43bd6795f5ec",
					Name:       "Book",
					Priority:   3,
					IsReserved: false,
				},
			},
		},
	)

	wishlist, items, err := svc.GetWishlistByToken(context.Background(), "token-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wishlist.ID == "" {
		t.Fatalf("wishlist expected to be returned")
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}
