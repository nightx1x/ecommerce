package service

import "errors"

var (
	//user errors
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserEmailRequired    = errors.New("user email is required")
	ErrUserPasswordRequired = errors.New("user password is required")
	ErrUserFNameRequired    = errors.New("user first name is required")
	ErrUserLNameRequired    = errors.New("user last name is required")

	//validation errors
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrInvalidRole        = errors.New("invalid role")
)
