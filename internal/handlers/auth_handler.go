package handlers

import (
	"errors"
	"net/http"

	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/middleware"
	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/service"
	"github.com/CodebyTecs/wishlist-service/pkg/httpx"
)

type AuthHandler struct {
	authService service.AuthService
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": domain.ErrInvalidRequest.Error()})
		return
	}

	accessToken, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, authResponse{AccessToken: accessToken})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": domain.ErrInvalidRequest.Error()})
		return
	}

	accessToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, authResponse{AccessToken: accessToken})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": domain.ErrUnauthorized.Error()})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"user_id": userID})
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRequest):
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrAlreadyExists):
		httpx.WriteJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrUnauthorized):
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
	default:
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}
