package http

import (
	"net/http"

	"github.com/CodebyTecs/wishlist-service/internal/handlers"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	wishlistHandler *handlers.WishlistHandler,
	requireAuth func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	if authHandler != nil {
		mux.HandleFunc("POST /auth/register", authHandler.Register)
		mux.HandleFunc("POST /auth/login", authHandler.Login)
	}

	if userHandler != nil && requireAuth != nil {
		mux.Handle("GET /users/me", requireAuth(http.HandlerFunc(userHandler.Me)))
	}
	if wishlistHandler != nil && requireAuth != nil {
		mux.Handle("POST /wishlists", requireAuth(http.HandlerFunc(wishlistHandler.Create)))
		mux.Handle("GET /wishlists", requireAuth(http.HandlerFunc(wishlistHandler.List)))
		mux.Handle("GET /wishlists/{id}", requireAuth(http.HandlerFunc(wishlistHandler.GetByID)))
		mux.Handle("PATCH /wishlists/{id}", requireAuth(http.HandlerFunc(wishlistHandler.UpdateByID)))
		mux.Handle("DELETE /wishlists/{id}", requireAuth(http.HandlerFunc(wishlistHandler.DeleteByID)))
	}

	return mux
}
