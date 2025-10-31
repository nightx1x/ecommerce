package service

import (
	"context"

	"github.com/google/uuid"
	repository "github.com/nightx1x/ecommerce/interval/repository/postgres"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, email, password string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) UserService {
	return &service{userRepo: userRepo}
}

// CreateUser implements UserService.
func (s *service) CreateUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	panic("unimplemented")
}

// DeleteUser implements UserService.
func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// GetUserByEmail implements UserService.
func (s *service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	panic("unimplemented")
}

// GetUserByID implements UserService.
func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	panic("unimplemented")
}

// UpdateUser implements UserService.
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, email, password string) error {
	panic("unimplemented")
}
