package response

import (
	"net/http"

	"github.com/CodebyTecs/wishlist-service/pkg/httpx"
)

type ErrorPayload struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, message string) {
	httpx.WriteJSON(w, status, ErrorPayload{Error: message})
}
