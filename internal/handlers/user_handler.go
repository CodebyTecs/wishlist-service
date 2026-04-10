package handlers

import (
	"errors"
	"net/http"

	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/middleware"
	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/service"
	"github.com/CodebyTecs/wishlist-service/pkg/httpx"
)

type UserHandler struct {
	users service.UserService
}

func NewUserHandler(users service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": domain.ErrUnauthorized.Error()})
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidRequest):
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			httpx.WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, user)
}
