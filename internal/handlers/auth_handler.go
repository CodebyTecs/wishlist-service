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

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	accessToken, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, dto.AuthResponse{AccessToken: accessToken})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	accessToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.AuthResponse{AccessToken: accessToken})
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRequest):
		response.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		response.WriteError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		response.WriteError(w, http.StatusUnauthorized, err.Error())
	default:
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
