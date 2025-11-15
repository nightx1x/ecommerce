package user

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserIDRequired = errors.New("user id is required")
	ErrInvalidEmail   = errors.New("invalid email format")
	ErrInvalidRole    = errors.New("invalid role value")

	ErrNoFields = errors.New("no fields to update")

	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPasswordRequired   = errors.New("password is required")
)
