package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/dto"
	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/middleware"
	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/response"
	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/service"
	"github.com/CodebyTecs/wishlist-service/pkg/httpx"
)

type WishlistItemHandler struct {
	items service.WishlistItemService
}

func NewWishlistItemHandler(items service.WishlistItemService) *WishlistItemHandler {
	return &WishlistItemHandler{items: items}
}

// Create godoc
// @Summary Create wishlist item
// @Description Создает позицию в указанном вишлисте текущего пользователя.
// @Tags WishlistItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param wishlistID path string true "Wishlist ID (UUID)"
// @Param request body dto.CreateWishlistItemRequest true "Wishlist item create payload"
// @Success 201 {object} domain.WishlistItem
// @Failure 400 {object} response.ErrorPayload
// @Failure 401 {object} response.ErrorPayload
// @Failure 404 {object} response.ErrorPayload
// @Failure 409 {object} response.ErrorPayload
// @Failure 500 {object} response.ErrorPayload
// @Router /wishlists/{wishlistID}/items [post]
func (h *WishlistItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	wishlistID := r.PathValue("wishlistID")

	var req dto.CreateWishlistItemRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	item, err := h.items.Create(r.Context(), userID, wishlistID, domain.CreateWishlistItemInput{
		Name:        req.Name,
		Description: req.Description,
		ProductLink: req.ProductLink,
		Priority:    req.Priority,
	})
	if err != nil {
		writeWishlistItemError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, item)
}

// List godoc
// @Summary List wishlist items
// @Description Возвращает позиции в указанном вишлисте текущего пользователя.
// @Tags WishlistItems
// @Produce json
// @Security BearerAuth
// @Param wishlistID path string true "Wishlist ID (UUID)"
// @Success 200 {array} domain.WishlistItem
// @Failure 400 {object} response.ErrorPayload
// @Failure 401 {object} response.ErrorPayload
// @Failure 500 {object} response.ErrorPayload
// @Router /wishlists/{wishlistID}/items [get]
func (h *WishlistItemHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	items, err := h.items.List(r.Context(), userID, r.PathValue("wishlistID"))
	if err != nil {
		writeWishlistItemError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, items)
}

// GetByID godoc
// @Summary Get wishlist item by ID
// @Description Возвращает позицию по ID из указанного вишлиста текущего пользователя.
// @Tags WishlistItems
// @Produce json
// @Security BearerAuth
// @Param wishlistID path string true "Wishlist ID (UUID)"
// @Param itemID path string true "Item ID (UUID)"
// @Success 200 {object} domain.WishlistItem
// @Failure 400 {object} response.ErrorPayload
// @Failure 401 {object} response.ErrorPayload
// @Failure 404 {object} response.ErrorPayload
// @Failure 500 {object} response.ErrorPayload
// @Router /wishlists/{wishlistID}/items/{itemID} [get]
func (h *WishlistItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	item, err := h.items.GetByID(r.Context(), userID, r.PathValue("wishlistID"), r.PathValue("itemID"))
	if err != nil {
		writeWishlistItemError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, item)
}

// UpdateByID godoc
// @Summary Update wishlist item
// @Description Частично обновляет позицию в указанном вишлисте текущего пользователя.
// @Tags WishlistItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param wishlistID path string true "Wishlist ID (UUID)"
// @Param itemID path string true "Item ID (UUID)"
// @Param request body dto.UpdateWishlistItemRequest true "Wishlist item patch payload"
// @Success 200 {object} domain.WishlistItem
// @Failure 400 {object} response.ErrorPayload
// @Failure 401 {object} response.ErrorPayload
// @Failure 404 {object} response.ErrorPayload
// @Failure 409 {object} response.ErrorPayload
// @Failure 500 {object} response.ErrorPayload
// @Router /wishlists/{wishlistID}/items/{itemID} [patch]
func (h *WishlistItemHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	var req dto.UpdateWishlistItemRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	trimmedName := strings.TrimSpace(req.Name)
	trimmedDescription := strings.TrimSpace(req.Description)
	trimmedProductLink := strings.TrimSpace(req.ProductLink)

	item, err := h.items.UpdateByID(
		r.Context(),
		userID,
		r.PathValue("wishlistID"),
		r.PathValue("itemID"),
		domain.UpdateWishlistItemInput{
			Name:              trimmedName,
			Description:       trimmedDescription,
			ProductLink:       trimmedProductLink,
			Priority:          req.Priority,
			UpdateName:        trimmedName != "",
			UpdateDescription: trimmedDescription != "",
			UpdateProductLink: trimmedProductLink != "",
			UpdatePriority:    req.Priority > 0,
		},
	)
	if err != nil {
		writeWishlistItemError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, item)
}

// DeleteByID godoc
// @Summary Delete wishlist item
// @Description Удаляет позицию из указанного вишлиста текущего пользователя.
// @Tags WishlistItems
// @Security BearerAuth
// @Param wishlistID path string true "Wishlist ID (UUID)"
// @Param itemID path string true "Item ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorPayload
// @Failure 401 {object} response.ErrorPayload
// @Failure 404 {object} response.ErrorPayload
// @Failure 500 {object} response.ErrorPayload
// @Router /wishlists/{wishlistID}/items/{itemID} [delete]
func (h *WishlistItemHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	if err := h.items.DeleteByID(r.Context(), userID, r.PathValue("wishlistID"), r.PathValue("itemID")); err != nil {
		writeWishlistItemError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeWishlistItemError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRequest):
		response.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		response.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrAlreadyExists), errors.Is(err, domain.ErrAlreadyReserved):
		response.WriteError(w, http.StatusConflict, err.Error())
	default:
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
