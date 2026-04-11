package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
)

type wishlistRepoStub struct {
	lastCreated  domain.Wishlist
	lastUpdate   domain.WishlistUpdate
	createErr    error
	listResult   []domain.Wishlist
	listErr      error
	getResult    domain.Wishlist
	getErr       error
	updateResult domain.Wishlist
	updateErr    error
	deleteErr    error
}

func (s *wishlistRepoStub) Create(_ context.Context, wishlist domain.Wishlist) (domain.Wishlist, error) {
	s.lastCreated = wishlist
	if s.createErr != nil {
		return domain.Wishlist{}, s.createErr
	}
	return wishlist, nil
}

func (s *wishlistRepoStub) ListByUserID(_ context.Context, _ string) ([]domain.Wishlist, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.listResult, nil
}

func (s *wishlistRepoStub) GetByIDAndUserID(_ context.Context, _, _ string) (domain.Wishlist, error) {
	if s.getErr != nil {
		return domain.Wishlist{}, s.getErr
	}
	return s.getResult, nil
}

func (s *wishlistRepoStub) GetByPublicToken(_ context.Context, _ string) (domain.Wishlist, error) {
	return domain.Wishlist{}, nil
}

func (s *wishlistRepoStub) UpdateByIDAndUserID(_ context.Context, _, _ string, update domain.WishlistUpdate) (domain.Wishlist, error) {
	s.lastUpdate = update
	if s.updateErr != nil {
		return domain.Wishlist{}, s.updateErr
	}
	return s.updateResult, nil
}

func (s *wishlistRepoStub) DeleteByIDAndUserID(_ context.Context, _, _ string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	return nil
}

func TestNewWishlistService(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{})
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
}

func TestWishlistServiceCreateInvalidUserID(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{})

	_, err := svc.Create(
		context.Background(),
		"bad-id",
		domain.CreateWishlistInput{Name: "Birthday", EventDate: time.Now().UTC()},
	)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistServiceCreateInvalidInput(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{})

	_, err := svc.Create(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		domain.CreateWishlistInput{Name: " ", EventDate: time.Time{}},
	)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistServiceCreateRepoError(t *testing.T) {
	errBoom := errors.New("create error")

	repo := &wishlistRepoStub{createErr: errBoom}
	svc := NewWishlistService(repo)
	_, err := svc.Create(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		domain.CreateWishlistInput{Name: "Birthday", EventDate: time.Now().UTC()},
	)
	if err != errBoom {
		t.Fatalf("expected create error, got: %v", err)
	}
}

func TestWishlistServiceCreateSuccess(t *testing.T) {
	repo := &wishlistRepoStub{}
	svc := NewWishlistService(repo)
	eventDate := time.Date(2026, 12, 10, 0, 0, 0, 0, time.UTC)

	wishlist, err := svc.Create(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		domain.CreateWishlistInput{Name: "  Birthday  ", Description: "  My gifts ", EventDate: eventDate},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wishlist.ID == "" || wishlist.PublicToken == "" {
		t.Fatalf("expected generated IDs, got: %+v", wishlist)
	}
	if repo.lastCreated.Name != "Birthday" {
		t.Fatalf("expected trimmed name, got: %q", repo.lastCreated.Name)
	}
	if repo.lastCreated.Description != "My gifts" {
		t.Fatalf("expected trimmed description, got: %q", repo.lastCreated.Description)
	}
}

func TestWishlistServiceListByUserIDInvalidRequest(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{})

	_, err := svc.ListByUserID(context.Background(), "bad")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistServiceListByUserIDSuccess(t *testing.T) {
	repo := &wishlistRepoStub{
		listResult: []domain.Wishlist{
			{ID: "65352d24-779d-45f0-a369-086c36e98241", Name: "Birthday"},
		},
	}
	svc := NewWishlistService(repo)

	wishlists, err := svc.ListByUserID(context.Background(), "9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(wishlists) != 1 {
		t.Fatalf("expected 1 wishlist, got %d", len(wishlists))
	}
}

func TestWishlistServiceGetByIDInvalidRequest(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{})

	_, err := svc.GetByID(context.Background(), "bad", "bad")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistServiceGetByIDSuccess(t *testing.T) {
	repo := &wishlistRepoStub{
		getResult: domain.Wishlist{ID: "65352d24-779d-45f0-a369-086c36e98241", Name: "Birthday"},
	}
	svc := NewWishlistService(repo)

	wishlist, err := svc.GetByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"65352d24-779d-45f0-a369-086c36e98241",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wishlist.ID == "" {
		t.Fatalf("expected wishlist")
	}
}

func TestWishlistServiceUpdateByIDInvalidRequest(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{})

	_, err := svc.UpdateByID(context.Background(), "bad", "bad", domain.UpdateWishlistInput{})
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistServiceUpdateByIDNoFields(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{})

	_, err := svc.UpdateByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"65352d24-779d-45f0-a369-086c36e98241",
		domain.UpdateWishlistInput{},
	)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistServiceUpdateByIDSuccess(t *testing.T) {
	repo := &wishlistRepoStub{
		updateResult: domain.Wishlist{
			ID:          "65352d24-779d-45f0-a369-086c36e98241",
			Name:        "Birthday Updated",
			Description: "Updated desc",
		},
	}
	svc := NewWishlistService(repo)

	wishlist, err := svc.UpdateByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"65352d24-779d-45f0-a369-086c36e98241",
		domain.UpdateWishlistInput{
			Name:              "  Birthday Updated ",
			Description:       " Updated desc ",
			UpdateName:        true,
			UpdateDescription: true,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wishlist.ID == "" {
		t.Fatalf("expected wishlist")
	}
	if repo.lastUpdate.Name != "Birthday Updated" {
		t.Fatalf("expected trimmed name in update payload, got: %q", repo.lastUpdate.Name)
	}
	if repo.lastUpdate.Description != "Updated desc" {
		t.Fatalf("expected trimmed description in update payload, got: %q", repo.lastUpdate.Description)
	}
}

func TestWishlistServiceDeleteByIDInvalidRequest(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{})

	err := svc.DeleteByID(context.Background(), "bad", "bad")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestWishlistServiceDeleteByIDPassThroughError(t *testing.T) {
	errBoom := errors.New("delete error")
	svc := NewWishlistService(&wishlistRepoStub{deleteErr: errBoom})

	err := svc.DeleteByID(
		context.Background(),
		"9a0ee5b9-e14a-47c0-b80a-f59de2dce5f7",
		"65352d24-779d-45f0-a369-086c36e98241",
	)
	if err != errBoom {
		t.Fatalf("expected delete error, got: %v", err)
	}
}
