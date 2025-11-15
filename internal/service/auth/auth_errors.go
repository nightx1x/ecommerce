package auth

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailRequired      = errors.New("user email is required")
	ErrPasswordRequired   = errors.New("user password is required")
	ErrWeakPassword       = errors.New("password weak")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrExpiredToken       = errors.New("token has expired")
)
