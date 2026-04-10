package domain

import "time"

type Wishlist struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	EventDate   time.Time `json:"event_date"`
	PublicToken string    `json:"public_token"`
	CreatedAt   time.Time `json:"created_at"`
}
