package service

import (
	"context"
	"errors"
	"testing"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
)

type wishlistItemRepoStub struct {
	lastCreated  domain.WishlistItem
	lastUpdate   domain.WishlistItemUpdate
	createErr    error
	listResult   []domain.WishlistItem
	listErr      error
	getResult    domain.WishlistItem
	getErr       error
	updateResult domain.WishlistItem
	updateErr    error
	deleteErr    error
}

func (s *wishlistItemRepoStub) Create(_ context.Context, _, _ string, item domain.WishlistItem) (domain.WishlistItem, error) {
	s.lastCreated = item
	if s.createErr != nil {
		return domain.WishlistItem{}, s.createErr
	}
	return item, nil
}

func (s *wishlistItemRepoStub) ListByWishlistAndUser(_ context.Context, _, _ string) ([]domain.WishlistItem, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.listResult, nil
}

func (s *wishlistItemRepoStub) GetByIDAndWishlistAndUser(_ context.Context, _, _, _ string) (domain.WishlistItem, error) {
	if s.getErr != nil {
		return domain.WishlistItem{}, s.getErr
	}
	return s.getResult, nil
}

func (s *wishlistItemRepoStub) UpdateByIDAndWishlistAndUser(_ context.Context, _, _, _ string, update domain.WishlistItemUpdate) (domain.WishlistItem, error) {
	s.lastUpdate = update
	if s.updateErr != nil {
		return domain.WishlistItem{}, s.updateErr
	}
	return s.updateResult, nil
}

func (s *wishlistItemRepoStub) DeleteByIDAndWishlistAndUser(_ context.Context, _, _, _ string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	return nil
}

func (s *wishlistItemRepoStub) ListByPublicToken(_ context.Context, _ string) ([]domain.WishlistItem, error) {
	return []domain.WishlistItem{}, nil
}

func (s *wishlistItemRepoStub) ReserveByPublicTokenAndItemID(_ context.Context, _, _ string) error {
	return nil
}

func TestNewWishlistItemService(t *testing.T) {
	svc := NewWishlistItemService(&wishlistItemRepoStub{})
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
}

func TestIsValidUserAndWishlistIDs(t *testing.T) {
	valid := isValidUserAndWishlistIDs(
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
	)
	if !valid {
		t.Fatalf("expected valid IDs")
	}

	invalid := isValidUserAndWishlistIDs("bad", "also-bad")
	if invalid {
		t.Fatalf("expected invalid IDs")
	}
}

func TestWishlistItemServiceCreateInvalidIDs(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	_, err := svc.Create(context.Background(), "bad", "bad", domain.CreateWishlistItemInput{Name: "Item", Priority: 3})
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceCreateInvalidName(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	_, err := svc.Create(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		domain.CreateWishlistItemInput{Name: "   ", Priority: 3},
	)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceCreateInvalidPriority(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	_, err := svc.Create(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		domain.CreateWishlistItemInput{
			Name:     "Keyboard",
			Priority: 10,
		},
	)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceCreateRepoError(t *testing.T) {
	errBoom := errors.New("create error")

	repo := &wishlistItemRepoStub{createErr: errBoom}
	svc := NewWishlistItemService(repo)

	_, err := svc.Create(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		domain.CreateWishlistItemInput{Name: "Keyboard", Priority: 5},
	)
	if err != errBoom {
		t.Fatalf("expected create error, got: %v", err)
	}
}

func TestWishlistItemServiceCreateSuccess(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	item, err := svc.Create(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		domain.CreateWishlistItemInput{
			Name:        "  Keyboard  ",
			Description: " 75% ",
			ProductLink: " https://example.com/item ",
			Priority:    5,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID == "" {
		t.Fatalf("expected generated item ID")
	}
	if item.Name != "Keyboard" {
		t.Fatalf("expected trimmed name, got: %q", item.Name)
	}
	if item.Description != "75%" {
		t.Fatalf("expected trimmed description, got: %q", item.Description)
	}
	if item.ProductLink != "https://example.com/item" {
		t.Fatalf("expected trimmed product link, got: %q", item.ProductLink)
	}
	if !item.CreatedAt.Equal(repo.lastCreated.CreatedAt) {
		t.Fatalf("expected item from service to be passed to repo")
	}
}

func TestWishlistItemServiceListInvalidRequest(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	_, err := svc.List(context.Background(), "bad", "bad")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceListSuccess(t *testing.T) {
	repo := &wishlistItemRepoStub{
		listResult: []domain.WishlistItem{
			{ID: "19f3638e-3274-4c99-8a34-43bd6795f5ec", Name: "Keyboard"},
		},
	}
	svc := NewWishlistItemService(repo)

	items, err := svc.List(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestWishlistItemServiceGetByIDInvalidRequest(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	_, err := svc.GetByID(context.Background(), "bad", "bad", "bad")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceGetByIDSuccess(t *testing.T) {
	repo := &wishlistItemRepoStub{
		getResult: domain.WishlistItem{
			ID:   "19f3638e-3274-4c99-8a34-43bd6795f5ec",
			Name: "Keyboard",
		},
	}
	svc := NewWishlistItemService(repo)

	item, err := svc.GetByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		"19f3638e-3274-4c99-8a34-43bd6795f5ec",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID == "" {
		t.Fatalf("expected item")
	}
}

func TestWishlistItemServiceUpdateInvalidIDs(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	_, err := svc.UpdateByID(context.Background(), "bad", "bad", "bad", domain.UpdateWishlistItemInput{})
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceUpdateWithoutFields(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	_, err := svc.UpdateByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		"19f3638e-3274-4c99-8a34-43bd6795f5ec",
		domain.UpdateWishlistItemInput{},
	)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceUpdateInvalidPriority(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	_, err := svc.UpdateByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		"19f3638e-3274-4c99-8a34-43bd6795f5ec",
		domain.UpdateWishlistItemInput{Priority: 10, UpdatePriority: true},
	)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceUpdateSuccess(t *testing.T) {
	repo := &wishlistItemRepoStub{
		updateResult: domain.WishlistItem{
			ID:          "19f3638e-3274-4c99-8a34-43bd6795f5ec",
			Description: "Updated item",
		},
	}
	svc := NewWishlistItemService(repo)

	item, err := svc.UpdateByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		"19f3638e-3274-4c99-8a34-43bd6795f5ec",
		domain.UpdateWishlistItemInput{
			Description:       " Updated item ",
			UpdateDescription: true,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID == "" {
		t.Fatalf("expected updated item")
	}
	if repo.lastUpdate.Description != "Updated item" {
		t.Fatalf("expected trimmed description in update payload, got: %q", repo.lastUpdate.Description)
	}
}

func TestWishlistItemServiceDeleteInvalidRequest(t *testing.T) {
	repo := &wishlistItemRepoStub{}
	svc := NewWishlistItemService(repo)

	err := svc.DeleteByID(context.Background(), "bad", "bad", "bad")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistItemServiceDeletePassThroughError(t *testing.T) {
	errBoom := errors.New("delete error")
	repo := &wishlistItemRepoStub{deleteErr: errBoom}
	svc := NewWishlistItemService(repo)

	err := svc.DeleteByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"6ef4f0a8-88ea-4ee8-911c-b5487f312380",
		"19f3638e-3274-4c99-8a34-43bd6795f5ec",
	)
	if err != errBoom {
		t.Fatalf("expected delete error, got: %v", err)
	}
}
