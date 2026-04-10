package domain

import "time"

type WishlistItem struct {
	ID          string     `json:"id"`
	WishlistID  string     `json:"wishlist_id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ProductLink string     `json:"product_link,omitempty"`
	Priority    int        `json:"priority"`
	IsReserved  bool       `json:"is_reserved"`
	ReservedAt  *time.Time `json:"reserved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
