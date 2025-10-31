package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	database "github.com/nightx1x/ecommerce/interval/db"
	models "github.com/nightx1x/ecommerce/interval/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	List(ctx context.Context, filter OrderFilter) ([]*models.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Cancel(ctx context.Context, id uuid.UUID) error
}
type orderRepo struct {
	db *database.DB
}

// Create implements OrderRepository.
func (o *orderRepo) Create(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (id, user_id, items, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := o.db.ExecContext(ctx, query,
		order.ID,
		order.UserID,
		order.Items,
		order.TotalAmount,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// Cancel implements OrderRepository.
func (o *orderRepo) Cancel(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE orders
		SET status = 'canceled', updated_at = NOW()
		WHERE id = $1
	`
	_, err := o.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}

// GetByID implements OrderRepository.
func (o *orderRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	query := `
		SELECT id, user_id, items, total_amount, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	row := o.db.QueryRowContext(ctx, query, id)
	var order models.Order
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Items,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}
	return &order, nil
}

// List implements OrderRepository.
func (o *orderRepo) List(ctx context.Context, filter OrderFilter) ([]*models.Order, error) {
	query := `
		SELECT id, user_id, items, total_amount, status, created_at, updated_at
		FROM orders
		WHERE ($1::uuid IS NULL OR user_id = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := o.db.QueryContext(ctx, query,
		filter.UserID,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()
	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Items,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

// UpdateStatus implements OrderRepository.
func (o *orderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := o.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
