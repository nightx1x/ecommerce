package cart

import "errors"

var (
	//Validate errors
	ErrUserIDRequired    = errors.New("user id is required")
	ErrItemIDRequired    = errors.New("item id is required")
	ErrProductIDRequired = errors.New("product id is required")
	ErrEmptyQuantity     = errors.New("quantity must be greater than 0")
	ErrInvalidQuantity   = errors.New("invalid quantity value")

	// Product error
	ErrInvalidProductID    = errors.New("invalid product id")
	ErrProductNotFound     = errors.New("product not found")
	ErrProductNotAvailable = errors.New("product not available")

	// Cart item error
	ErrItemNotFound = errors.New("cart item not found")
)
