package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	models "github.com/nightx1x/ecommerce/interval/domain"
	repository "github.com/nightx1x/ecommerce/interval/repository/postgres"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	ListProducts(ctx context.Context, filter ProductFilter) (*ProductListResponse, error)
	SearchProducts(ctx context.Context, query string, limit, offset int) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*models.Product, error)
	CheckAvailability(ctx context.Context, id uuid.UUID, quantity int) (bool, error)
	ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error
	ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error
}

type CreateProductRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	Price       float64    `json:"price" validate:"required,gt=0"`
	Stock       int        `json:"stock" validate:"required,gte=0"`
	CategoryID  *uuid.UUID `json:"category_id" validate:"omitempty,uuid"`
	ImageURL    string     `json:"image_url" validate:"omitempty,url"`
}

// UpdateProductRequest is the DTO for updating a product
type UpdateProductRequest struct {
	Name        *string    `json:"name" validate:"omitempty,min=3,max=255"`
	Description *string    `json:"description" validate:"omitempty,max=1000"`
	Price       *float64   `json:"price" validate:"omitempty,gt=0"`
	Stock       *int       `json:"stock" validate:"omitempty,gte=0"`
	CategoryID  *uuid.UUID `json:"category_id" validate:"omitempty,uuid"`
	ImageURL    *string    `json:"image_url" validate:"omitempty,url"`
}

// ProductFilter is the DTO for filtering products
type ProductFilter struct {
	CategoryID *uuid.UUID `json:"category_id"`
	MinPrice   *float64   `json:"min_price" validate:"omitempty,gte=0"`
	MaxPrice   *float64   `json:"max_price" validate:"omitempty,gte=0"`
	Search     string     `json:"search"`
	InStock    *bool      `json:"in_stock"`
	OrderBy    string     `json:"order_by" validate:"omitempty,oneof=price_asc price_desc name_asc name_desc created_at_asc created_at_desc"`
	Limit      int        `json:"limit" validate:"required,min=1,max=100"`
	Offset     int        `json:"offset" validate:"gte=0"`
}

// ProductListResponse contains paginated products and metadata
type ProductListResponse struct {
	Products []*models.Product `json:"products"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

type service struct {
	productRepo repository.ProductRepository
}

// CheckAvailability implements ProductService.
func (s *service) CheckAvailability(ctx context.Context, id uuid.UUID, quantity int) (bool, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return false, ErrProductNotFound
	}
	if product.Stock >= quantity {
		return true, nil
	}
	return product.Stock < quantity, nil
}

// GetProductByID implements ProductService.
func (s *service) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

// ListProducts implements ProductService.
func (s *service) ListProducts(ctx context.Context, filter ProductFilter) (*ProductListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit >= 100 {
		filter.Limit = 100
	}

	if filter.Offset > 0 {
		filter.Limit = 20
	}

	if filter.MinPrice != nil && *filter.MinPrice < 0 {
		return nil, ErrInvalidPrice
	}
	if filter.MaxPrice != nil && *filter.MaxPrice < 0 {
		return nil, ErrInvalidPrice
	}

	repoFilter := models.ListFilter{
		CategoryID: filter.CategoryID,
		MinPrice:   filter.MinPrice,
		MaxPrice:   filter.MaxPrice,
		Search:     filter.Search,
		OrderBy:    filter.OrderBy,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	}

	product, err := s.productRepo.List(ctx, &repoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	if filter.InStock != nil {
		filteredProducts := make([]*models.Product, 0)
		for _, p := range product {
			if p.Stock > 0 {
				filteredProducts = append(filteredProducts, p)
			}
		}
		product = filteredProducts
	}
	return &ProductListResponse{
		Products: product,
		Total:    len(product),
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil
}

// ReleaseStock implements ProductService.
func (s *service) ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	product.Stock += quantity
	return s.productRepo.Update(ctx, product)

}

// ReserveStock implements ProductService.
func (s *service) ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product.Stock < quantity {
		return ErrInvalidStock
	}
	product.Stock -= quantity
	return s.productRepo.Update(ctx, product)

}

// SearchProducts implements ProductService.
func (s *service) SearchProducts(ctx context.Context, query string, limit int, offset int) ([]*models.Product, error) {
	if query == "" {
		return nil, ErrProductNameRequired
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	filter := models.ListFilter{
		Search: query,
		Limit:  limit,
		Offset: offset,
	}
	products, err := s.productRepo.List(ctx, &filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	return products, nil
}

// UpdateProduct implements ProductService.
func (s *service) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, ErrInvalidPrice
		}
		product.Price = *req.Price
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, ErrInvalidStock
		}
		product.Stock = *req.Stock
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.ImageURL != nil {
		product.ImageURL = req.ImageURL
	}
	err = s.productRepo.Update(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return product, nil
}

func NewService(productRepo repository.ProductRepository) ProductService {
	return &service{
		productRepo: productRepo,
	}
}

func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error) {
	if req.Name == "" || len(req.Name) < 3 {
		return nil, ErrInvalidName
	}
	if req.Price <= 0 {
		return nil, ErrInvalidPrice
	}
	if req.Stock < 0 {
		return nil, ErrInvalidStock
	}
	var description *string
	if req.Description != "" {
		description = &req.Description
	}
	var imageURL *string
	if req.ImageURL != "" {
		imageURL = &req.ImageURL
	}

	product := &models.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		ImageURL:    imageURL,
	}
	err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}
