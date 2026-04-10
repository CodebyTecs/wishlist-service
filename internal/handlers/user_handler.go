package handlers

import (
	"errors"
	"net/http"

	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/middleware"
	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/response"
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
		response.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidRequest):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrNotFound):
			response.WriteError(w, http.StatusNotFound, err.Error())
		default:
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, user)
}
