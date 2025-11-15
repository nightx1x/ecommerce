package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	models "github.com/nightx1x/ecommerce/internal/domain"
	repository "github.com/nightx1x/ecommerce/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest, isAdmin bool) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UpdateUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6,max=100"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Role      *string `json:"role" validate:"omitempty,oneof=user admin"`
}
type UpdateOwnUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6,max=100"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=100"`
}

type UserFilter struct {
	Search  string  `json:"search"`
	Role    *string `json:"role" validate:"omitempty,oneof=user admin"`
	Limit   int     `json:"limit" validate:"required,min=1,max=100"`
	Offset  int     `json:"offset" validate:"gte=0"`
	OrderBy string  `json:"order_by" validate:"omitempty,oneof=created_at_asc created_at_desc name_asc name_desc"`
}
type UserListResponse struct {
	Users  []*models.User `json:"users"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) UserService {
	return &service{userRepo: userRepo}
}

// DeleteUser
func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {

	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delte user: %w", err)
	}

	return nil
}

// GetUser
func (s *service) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {

	if id == uuid.Nil {
		return nil, ErrUserNotFound
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("faield to get user: %w", err)
	}

	return user, nil
}

// ListUser
func (s *service) ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error) {

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	users, err := s.userRepo.List(ctx, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	filtered := []*models.User{}

	for _, u := range users {
		if filter.Role != nil && u.Role != *filter.Role {
			continue
		}
		if filter.Search != "" && !strings.Contains(strings.ToLower(u.FirstName),
			strings.ToLower(filter.Search)) && !strings.Contains(strings.ToLower(u.LastName), strings.ToLower(filter.Search)) {
			continue
		}
		filtered = append(filtered, u)
	}

	resp := &UserListResponse{
		Users: filtered,
		Total: len(filtered),
	}

	return resp, nil
}

// UpdateUser
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest, isAdmin bool) (*models.User, error) {

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Email == nil && req.FirstName == nil && req.LastName == nil &&
		req.Password == nil && req.Role == nil {
		return nil, ErrNoFields
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	if req.Role != nil && isAdmin {
		user.Role = *req.Role
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
