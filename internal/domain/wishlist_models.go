package domain

import "time"

type CreateWishlistInput struct {
	Name        string
	Description string
	EventDate   time.Time
}

type UpdateWishlistInput struct {
	Name              string
	Description       string
	EventDate         time.Time
	UpdateName        bool
	UpdateDescription bool
	UpdateEventDate   bool
}

type WishlistUpdate struct {
	Name              string
	Description       string
	EventDate         time.Time
	UpdateName        bool
	UpdateDescription bool
	UpdateEventDate   bool
}
