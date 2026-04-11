package http

import (
	"net/http"

	"github.com/CodebyTecs/wishlist-service/internal/handlers"
)

func NewRouter(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, wishlistHandler *handlers.WishlistHandler, wishlistItemHandler *handlers.WishlistItemHandler, publicHandler *handlers.PublicHandler, requireAuth func(http.Handler) http.Handler) http.Handler {
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
	if wishlistItemHandler != nil && requireAuth != nil {
		mux.Handle("POST /wishlists/{wishlistID}/items", requireAuth(http.HandlerFunc(wishlistItemHandler.Create)))
		mux.Handle("GET /wishlists/{wishlistID}/items", requireAuth(http.HandlerFunc(wishlistItemHandler.List)))
		mux.Handle("GET /wishlists/{wishlistID}/items/{itemID}", requireAuth(http.HandlerFunc(wishlistItemHandler.GetByID)))
		mux.Handle("PATCH /wishlists/{wishlistID}/items/{itemID}", requireAuth(http.HandlerFunc(wishlistItemHandler.UpdateByID)))
		mux.Handle("DELETE /wishlists/{wishlistID}/items/{itemID}", requireAuth(http.HandlerFunc(wishlistItemHandler.DeleteByID)))
	}
	if publicHandler != nil {
		mux.HandleFunc("GET /public/{token}", publicHandler.GetWishlistByToken)
		mux.HandleFunc("POST /public/{token}/reserve/{itemID}", publicHandler.ReserveByTokenAndItemID)
	}

	return mux
}
