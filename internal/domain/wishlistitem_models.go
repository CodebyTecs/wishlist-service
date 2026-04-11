package domain

type CreateWishlistItemInput struct {
	Name        string
	Description string
	ProductLink string
	Priority    int
}

type UpdateWishlistItemInput struct {
	Name              string
	Description       string
	ProductLink       string
	Priority          int
	UpdateName        bool
	UpdateDescription bool
	UpdateProductLink bool
	UpdatePriority    bool
}

type WishlistItemUpdate struct {
	Name              string
	Description       string
	ProductLink       string
	Priority          int
	UpdateName        bool
	UpdateDescription bool
	UpdateProductLink bool
	UpdatePriority    bool
}
