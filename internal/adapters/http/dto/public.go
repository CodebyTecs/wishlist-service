package dto

import "time"

type PublicWishlistItem struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ProductLink string     `json:"product_link,omitempty"`
	Priority    int        `json:"priority"`
	IsReserved  bool       `json:"is_reserved"`
	ReservedAt  *time.Time `json:"reserved_at,omitempty"`
}

type PublicWishlistResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	EventDate   time.Time            `json:"event_date"`
	Items       []PublicWishlistItem `json:"items"`
}
