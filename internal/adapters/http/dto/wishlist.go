package dto

type CreateWishlistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
}

type UpdateWishlistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
}
