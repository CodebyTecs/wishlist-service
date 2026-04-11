package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/dto"
	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/middleware"
	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/response"
	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/service"
	"github.com/CodebyTecs/wishlist-service/pkg/httpx"
)

const dateLayout = "2006-01-02"

type WishlistHandler struct {
	wishlists service.WishlistService
}

func NewWishlistHandler(wishlists service.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlists: wishlists}
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	var req dto.CreateWishlistRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	eventDate, err := time.Parse(dateLayout, req.EventDate)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "event_date must be in YYYY-MM-DD format")
		return
	}

	wishlist, err := h.wishlists.Create(r.Context(), userID, domain.CreateWishlistInput{
		Name:        req.Name,
		Description: req.Description,
		EventDate:   eventDate,
	})
	if err != nil {
		writeWishlistError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, wishlist)
}

func (h *WishlistHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	wishlists, err := h.wishlists.ListByUserID(r.Context(), userID)
	if err != nil {
		writeWishlistError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, wishlists)
}

func (h *WishlistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	wishlistID := r.PathValue("id")
	wishlist, err := h.wishlists.GetByID(r.Context(), userID, wishlistID)
	if err != nil {
		writeWishlistError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, wishlist)
}

func (h *WishlistHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	var req dto.UpdateWishlistRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	trimmedName := strings.TrimSpace(req.Name)
	trimmedDescription := strings.TrimSpace(req.Description)

	updateName := trimmedName != ""
	updateDescription := trimmedDescription != ""
	updateEventDate := strings.TrimSpace(req.EventDate) != ""

	var parsedDate time.Time
	if updateEventDate {
		d, err := time.Parse(dateLayout, req.EventDate)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "event_date must be in YYYY-MM-DD format")
			return
		}
		parsedDate = d
	}

	wishlist, err := h.wishlists.UpdateByID(r.Context(), userID, r.PathValue("id"), domain.UpdateWishlistInput{
		Name:              trimmedName,
		Description:       trimmedDescription,
		EventDate:         parsedDate,
		UpdateName:        updateName,
		UpdateDescription: updateDescription,
		UpdateEventDate:   updateEventDate,
	})
	if err != nil {
		writeWishlistError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, wishlist)
}

func (h *WishlistHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	if err := h.wishlists.DeleteByID(r.Context(), userID, r.PathValue("id")); err != nil {
		writeWishlistError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeWishlistError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRequest):
		response.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		response.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrAlreadyExists):
		response.WriteError(w, http.StatusConflict, err.Error())
	default:
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
