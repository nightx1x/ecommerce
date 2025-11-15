package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	productSrv "github.com/nightx1x/ecommerce/internal/service/product"
)

type ProductHandler struct {
	ProductSrv productSrv.ProductService
}

func NewProductHandler(srv productSrv.ProductService) *ProductHandler {
	return &ProductHandler{ProductSrv: srv}
}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/products", h.ListProducts)
		r.Get("/products/search", h.SearchProduct)
		r.Get("/products/{id}", h.GetProduct)
		r.Get("/categories", h.ListCategories)
	})

	r.Group(func(r chi.Router) {
		r.Post("/admin/products", h.CreateProduct)
		r.Put("/admin/products/{id}", h.UpdateProduct)
	})
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Returns detailed product information by its unique identifier
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} models.Product
// @Failure 400 {object} ErrorResponse "Invalid product ID"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /products/{id} [get]

// handler

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}

	product, err := h.ProductSrv.GetProductByID(r.Context(), id)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, product)
}

// ListProducts godoc
// @Summary Get list of products
// @Description Returns a list of products with optional filtering by category, price, availability, and pagination
// @Tags products
// @Accept json
// @Produce json
// @Param category_id query string false "Category ID (UUID)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param search query string false "Search query"
// @Param in_stock query boolean false "Only products in stock"
// @Param limit query integer false "Number of items per page" default(20) minimum(1) maximum(100)
// @Param offset query integer false "Offset for pagination" default(0) minimum(0)
// @Success 200 {object} product.ProductListResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	filter := productSrv.ProductFilter{
		Limit:  20,
		Offset: 0,
	}

	// categoryID
	if categoryIDStr := r.URL.Query().Get("category_id"); categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid category ID")
		}
		filter.CategoryID = &categoryID
	}

	// minPrice
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid min_price")
		}
		filter.MinPrice = &minPrice
	}

	// maxPrice
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid max_price")
		}
		filter.MaxPrice = &maxPrice
	}

	// Search
	filter.Search = r.URL.Query().Get("search")

	// InStock
	if inStockStr := r.URL.Query().Get("in_stock"); inStockStr != "" {
		instock := inStockStr == "true"
		filter.InStock = &instock
	}

	// Limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(w, http.StatusBadRequest, "Invalid limit")
		}
		filter.Limit = limit
	}

	// Offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "Invalid offset")
		}
		filter.Offset = offset
	}

	response, err := h.ProductSrv.ListProducts(r.Context(), filter)
	if err != nil {
		handlerServiceError(w, err)
	}
	respondJSON(w, http.StatusOK, response)
}

// ListCategories godoc
// @Summary Get list of categories
// @Description Returns a list of all product categories
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} interface{}
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Router /categories [get]
func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories := []interface{}{}
	respondJSON(w, http.StatusOK, categories)
}

// / SearchProduct godoc
// @Summary Search products
// @Description Returns a list of products that match the search query
// @Tags products
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query integer false "Number of items per page" default(20) minimum(1) maximum(100)
// @Param offset query integer false "Offset for pagination" default(0) minimum(0)
// @Success 200 {array} models.Product
// @Failure 400 {object} http.ErrorResponse "Search query is required or invalid parameters"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Router /products/search [get]
func (h *ProductHandler) SearchProduct(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 {
			respondError(w, http.StatusBadRequest, "Invalid limit")
		}
		limit = l
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			respondError(w, http.StatusBadRequest, "Invalid offset")
		}
		offset = o
	}

	products, err := h.ProductSrv.SearchProducts(r.Context(), query, limit, offset)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, products)
}

// ADMIN PART
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req productSrv.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createProd, err := h.ProductSrv.CreateProduct(r.Context(), req)
	if err != nil {
		handlerServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, createProd)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req productSrv.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updProd, err := h.ProductSrv.UpdateProduct(r.Context(), id, req)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updProd)
}

// utils
func handlerServiceError(w http.ResponseWriter, err error) {
	switch err {
	case productSrv.ErrProductNotFound:
		respondError(w, http.StatusNotFound, err.Error())
	case productSrv.ErrProductNameRequired,
		productSrv.ErrInvalidPrice,
		productSrv.ErrInvalidStock,
		productSrv.ErrInvalidQuantity:
		respondError(w, http.StatusBadRequest, err.Error())
	case productSrv.ErrInsufficientStock:
		respondError(w, http.StatusConflict, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}
