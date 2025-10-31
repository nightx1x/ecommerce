package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	database "github.com/nightx1x/ecommerce/interval/db"
	models "github.com/nightx1x/ecommerce/interval/domain"
)

type CartRepository interface {
	CreateItem(ctx context.Context, cart *models.Cart) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error)
	UpdateItem(ctx context.Context, cart *models.Cart) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
}
type cartRepo struct {
	db *database.DB
}

// Create implements CartRepository.
func (c *cartRepo) CreateItem(ctx context.Context, cart *models.Cart) error {
	query := `
		INSERT INTO carts (id, user_id, items, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := c.db.ExecContext(ctx, query,
		cart.ID,
		cart.UserID,
		cart.Items,
		cart.CreatedAt,
		cart.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create cart: %w", err)
	}
	return nil
}

// Delete implements CartRepository.
func (c *cartRepo) DeleteItem(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM carts
		WHERE id = $1
	`
	_, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}
	return nil
}

// GetByID implements CartRepository.
func (c *cartRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error) {
	query := `
		SELECT id, user_id, items, created_at, updated_at
		FROM carts
		WHERE id = $1
	`
	row := c.db.QueryRowContext(ctx, query, id)
	cart := &models.Cart{}
	err := row.Scan(
		&cart.ID,
		&cart.UserID,
		&cart.Items,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	return cart, nil
}

// Update implements CartRepository.
func (c *cartRepo) UpdateItem(ctx context.Context, cart *models.Cart) error {
	query := `
		UPDATE carts
		SET user_id = $1, items = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := c.db.ExecContext(ctx, query,
		cart.UserID,
		cart.Items,
		cart.UpdatedAt,
		cart.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}
	return nil
}

func NewCartRepository(db *database.DB) CartRepository {
	return &cartRepo{db: db}
}
