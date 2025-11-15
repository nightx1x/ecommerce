package product

import (
	"context"
	"testing"

	"github.com/google/uuid"
	models "github.com/nightx1x/ecommerce/interval/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) List(ctx context.Context, filter *models.ListFilter) ([]*models.Product, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateProduct_123(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	req := CreateProductRequest{
		Name:        "Test Product",
		Description: "Test description",
		Price:       99.99,
		Stock:       10,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Product")).Return(nil)
	//Act
	product, err := service.CreateProduct(ctx, req)

	//Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Stock, product.Stock)
	assert.Equal(t, req.Price, product.Price)
	mockRepo.AssertExpectations(t)

}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	req := CreateProductRequest{
		Name:        "Test Product",
		Description: "Test description",
		Price:       99.99,
		Stock:       10,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Product")).Return(nil)
	//Act
	product, err := service.CreateProduct(ctx, req)

	//Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Stock, product.Stock)
	assert.Equal(t, req.Price, product.Price)
	mockRepo.AssertExpectations(t)

}

func TestCreateProduct_InvalidPrice(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	req := CreateProductRequest{
		Name:  "",
		Price: -99.99,
		Stock: 10,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Product")).Return(nil)
	//Act
	product, err := service.CreateProduct(ctx, req)

	//Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, ErrInvalidPrice, err)
	mockRepo.AssertNotCalled(t, "Create")

}

func TestCreateProduct_NotAvailable(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	productId := uuid.New()
	product := &models.Product{
		ID:    productId,
		Name:  "Test Product",
		Stock: 3,
	}

	mockRepo.On("GetByID", ctx, productId).Return(product, nil)
	//Act
	available, err := service.CheckAvailability(ctx, productId, 5)

	//Assert
	assert.NoError(t, err)
	assert.False(t, available)

	mockRepo.AssertExpectations(t)

}
func TestCreateProduct_NotQuantity(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	productId := uuid.New()

	//Act
	available, err := service.CheckAvailability(ctx, productId, 0)

	//Assert
	assert.Error(t, err)
	assert.False(t, available)
	assert.Equal(t, ErrInvalidQuantity, err)

	mockRepo.AssertNotCalled(t, "GetByID")
}
func TestCreateProduct_NoеQuantity(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	productId := uuid.New()

	//Act
	available, err := service.CheckAvailability(ctx, productId, 0)

	//Assert
	assert.Error(t, err)
	assert.False(t, available)
	assert.Equal(t, ErrInvalidQuantity, err)

	mockRepo.AssertNotCalled(t, "GetByID")
}
