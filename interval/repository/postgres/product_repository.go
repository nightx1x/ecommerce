package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	database "github.com/nightx1x/ecommerce/interval/db"
	models "github.com/nightx1x/ecommerce/interval/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	List(ctx context.Context, filter *models.ListFilter) ([]*models.Product, error)
	Update(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productRepo struct {
	db *database.DB
}

// Create implements ProductRepository.
func (p *productRepo) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, stock, category_id, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := p.db.ExecContext(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.ImageURL,
		product.CreatedAt,
		product.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// Delete implements ProductRepository.
func (p *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM products
		WHERE id = $1
	`
	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// GetByID implements ProductRepository.
func (p *productRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	query := `
        SELECT id, name, description, price, stock, category_id, image_url, created_at, updated_at
        FROM products
        WHERE id = $1
    `

	err := p.db.GetContext(ctx, &product, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// List implements ProductRepository.
func (p *productRepo) List(ctx context.Context, filter *models.ListFilter) ([]*models.Product, error) {
	query := `
	SELECT id, name, description, price, stock, category_id, image_url, created_at, updated_at
		FROM products
		WHERE 1=1
	`

	args := []interface{}{}
	argsCount := 1

	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argsCount)
		args = append(args, *filter.CategoryID)
		argsCount++
	}

	if filter.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argsCount)
		args = append(args, *filter.MinPrice)
		argsCount++
	}

	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argsCount)
		args = append(args, *filter.MaxPrice)
		argsCount++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d OR description ILIKE $%d", argsCount, argsCount)
		args = append(args, "%"+filter.Search+"%")
		argsCount++
	}

	orderBy := "created_at DESC"
	if filter.OrderBy != "" {
		allowedOrders := map[string]bool{
			"price_asc":  true,
			"price_desc": true,
			"name_asc":   true,
			"name_desc":  true,
		}
		if allowedOrders[filter.OrderBy] {
			parts := strings.Split(filter.OrderBy, "_")
			orderBy = parts[0] + " " + strings.ToUpper(parts[1])
		}
	}
	query += " ORDER BY " + orderBy
	// Pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argsCount)
		args = append(args, filter.Limit)
		argsCount++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argsCount)
		args = append(args, filter.Offset)
		argsCount++
	}

	var products []*models.Product
	err := p.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}

// Update implements ProductRepository.
func (p *productRepo) Update(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE products
		SET updated_at = NOW()
		WHERE id = $1
	`
	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// UpdateStock implements ProductRepository.
func (p *productRepo) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock + $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := p.db.ExecContext(ctx, query, quantity, id)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}
	return nil
}

func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepo{db: db}
}
