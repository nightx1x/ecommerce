package service

import "errors"

// Domain errors for product service
var (
	// Product errors
	ErrProductNotFound     = errors.New("product not found")
	ErrProductNameRequired = errors.New("product name is required")
	ErrProductNotAvailable = errors.New("product is not available")

	// Validation errors
	ErrInvalidPrice    = errors.New("price must be greater than 0")
	ErrInvalidStock    = errors.New("stock must be non-negative")
	ErrInvalidQuantity = errors.New("quantity must be greater than 0")

	// Stock errors
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidName       = errors.New("invalid product name")
)
