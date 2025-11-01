package order

import (
	"context"

	"github.com/google/uuid"
	models "github.com/nightx1x/ecommerce/interval/domain"
	repository "github.com/nightx1x/ecommerce/interval/repository/postgres"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderRequset) (*OrderResponse, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*OrderResponse, error)
	ListOrder(ctx context.Context, filter *models.OrderFilter) ([]*OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) (*OrderResponse, error)
	CancelOrder(ctx context.Context, id uuid.UUID) error
}

type service struct {
	orderRepo repository.OrderRepository
}

func NewService(orderRepo repository.OrderRepository) OrderService {
	return &service{orderRepo: orderRepo}
}

// CancelOrder implements OrderService.
func (s *service) CancelOrder(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// CreateOrder implements OrderService.
func (s *service) CreateOrder(ctx context.Context, req CreateOrderRequset) (*OrderResponse, error) {
	panic("unimplemented")
}

// GetOrder implements OrderService.
func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (*OrderResponse, error) {
	panic("unimplemented")
}

// ListOrder implements OrderService.
func (s *service) ListOrder(ctx context.Context, filter *models.OrderFilter) ([]*OrderResponse, error) {
	panic("unimplemented")
}

// UpdateOrderStatus implements OrderService.
func (s *service) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) (*OrderResponse, error) {
	panic("unimplemented")
}
