package dto

type CreateWishlistItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProductLink string `json:"product_link"`
	Priority    int    `json:"priority"`
}

type UpdateWishlistItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProductLink string `json:"product_link"`
	Priority    int    `json:"priority"`
}
