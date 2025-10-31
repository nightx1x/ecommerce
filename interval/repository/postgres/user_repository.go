package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	database "github.com/nightx1x/ecommerce/interval/db"
	models "github.com/nightx1x/ecommerce/interval/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
type userRepo struct {
	db *database.DB
}

// Create implements UserRepository.
func (u *userRepo) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := u.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Delete implements UserRepository.
func (u *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetByID implements UserRepository.
func (u *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := u.db.QueryRowContext(ctx, query, id)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetByEmail implements UserRepository.
func (u *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	row := u.db.QueryRowContext(ctx, query, email)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// Update implements UserRepository.
func (u *userRepo) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := u.db.ExecContext(ctx, query,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
func NewUserRepository(db *database.DB) UserRepository {
	return &userRepo{db: db}
}
