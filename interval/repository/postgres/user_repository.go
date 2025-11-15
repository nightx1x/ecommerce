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
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}

type UserRepo struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &UserRepo{db: db}
}

// Create створює нового користувача
func (u *UserRepo) Create(ctx context.Context, user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `
	INSERT INTO users (id, email,password_hash, first_name, last_name, role, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())
	`

	_, err := u.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Delete видаляє користувача за його ID
func (u *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM users WHERE id = $1
	`

	res, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("users not found")
	}

	return nil
}

// GetByEmail повертає користувача за його email
func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, email,password_hash, first_name, last_name, role, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	err := u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by email: %w", err)
	}

	return &user, nil
}

// GetByID повертає користувача за його ID
func (u *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, email,password_hash, first_name, last_name, role, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	err := u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by id: %w", err)
	}

	return &user, nil
}

// List повертає список користувачів
func (u *UserRepo) List(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	query := `
	SELECT id, email,password_hash, first_name, last_name, role, created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	var users []*models.User
	err := u.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Update оновлює дані користувача
func (u *UserRepo) Update(ctx context.Context, user *models.User) error {
	query := `
	UPDATE users
	SET email = $1,
		password_hash = $2,
		first_name = $3,
		last_name = $4,
		role = $5,
		updated_at = NOW()
	WHERE id = $6	
	`
	res, err := u.db.ExecContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("users not found")
	}
	return nil
}
