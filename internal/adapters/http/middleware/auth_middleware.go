package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/CodebyTecs/wishlist-service/internal/service"
	"github.com/CodebyTecs/wishlist-service/pkg/httpx"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

type AuthMiddleware struct {
	tokens service.TokenService
}

func NewAuthMiddleware(tokens service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokens: tokens}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := extractBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": domain.ErrUnauthorized.Error()})
			return
		}

		userID, err := m.tokens.ParseAccessToken(tokenString)
		if err != nil {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": domain.ErrUnauthorized.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}

func extractBearerToken(authorization string) (string, error) {
	parts := strings.SplitN(strings.TrimSpace(authorization), " ", 2)
	if len(parts) != 2 {
		return "", domain.ErrUnauthorized
	}
	if !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", domain.ErrUnauthorized
	}
	return strings.TrimSpace(parts[1]), nil
}
