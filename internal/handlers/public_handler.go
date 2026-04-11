package handlers

import (
	"errors"
	"net/http"

	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/dto"
	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/response"
	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/service"
	"github.com/CodebyTecs/wishlist-service/pkg/httpx"
)

type PublicHandler struct {
	public service.PublicService
}

func NewPublicHandler(public service.PublicService) *PublicHandler {
	return &PublicHandler{public: public}
}

func (h *PublicHandler) GetWishlistByToken(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	wishlist, items, err := h.public.GetWishlistByToken(r.Context(), token)
	if err != nil {
		writePublicError(w, err)
		return
	}

	publicItems := make([]dto.PublicWishlistItem, 0, len(items))
	for _, item := range items {
		publicItems = append(publicItems, dto.PublicWishlistItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			ProductLink: item.ProductLink,
			Priority:    item.Priority,
			IsReserved:  item.IsReserved,
			ReservedAt:  item.ReservedAt,
		})
	}

	httpx.WriteJSON(w, http.StatusOK, dto.PublicWishlistResponse{
		ID:          wishlist.ID,
		Name:        wishlist.Name,
		Description: wishlist.Description,
		EventDate:   wishlist.EventDate,
		Items:       publicItems,
	})
}

func (h *PublicHandler) ReserveByTokenAndItemID(w http.ResponseWriter, r *http.Request) {
	err := h.public.ReserveItem(r.Context(), r.PathValue("token"), r.PathValue("itemID"))
	if err != nil {
		writePublicError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "reserved"})
}

func writePublicError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRequest):
		response.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		response.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyReserved):
		response.WriteError(w, http.StatusConflict, err.Error())
	default:
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
