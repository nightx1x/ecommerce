package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	orderSrv "github.com/nightx1x/ecommerce/internal/service/order"
	productSrv "github.com/nightx1x/ecommerce/internal/service/product"
	userSrv "github.com/nightx1x/ecommerce/internal/service/user"
)

type AdminHandler struct {
	productSrv productSrv.ProductService
	orderSrv   orderSrv.OrderService
	userSrv    userSrv.UserService
}

func NewAdminHandler(productSrv productSrv.ProductService, orderSrv orderSrv.OrderService, userSrv userSrv.UserService) *AdminHandler {
	return &AdminHandler{productSrv: productSrv,
		orderSrv: orderSrv,
		userSrv:  userSrv}
}

func (h *AdminHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		//Product
		r.Post("/products", h.CreateProduct)
		r.Put("/products/{id}", h.UpdateProduct)
		//r.Delete("/products/{id}", h.DeleteProduct) // Uncomment and implement DeleteProduct if needed

		//Order
		r.Get("/orders", h.ListAllOrder)
		r.Put("/orders/{id}/status", h.UpdateOrderStatus)

		//User
		r.Get("/users/{id}", h.GetUser)
		r.Get("/users", h.ListUsers)
		r.Delete("/users/{id}", h.DeleteUser)
		r.Put("/users/{id}", h.UpdateUser)
	})
}

// Products

// CreateProduct godoc
// @Summary Product Creation (Admin)
// @Description Creates a new product
// @Tags admin
// @Accept json
// @Produce json
// @Param product body product.CreateProductRequest true "Product data"
// @Success 200 {object} models.Product
// @Failure 400 {object} http.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /admin/products [post]
func (h *AdminHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req productSrv.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createProd, err := h.productSrv.CreateProduct(r.Context(), req)
	if err != nil {
		handlerServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, createProd)
}

// UpdateProduct godoc
// @Summary Product Update (Admin)
// @Description Updates an existing product
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body product.UpdateProductRequest true "Data for update"
// @Success 200 {object} models.Product
// @Failure 400 {object} http.ErrorResponse "Invalid request body or ID"
// @Failure 404 {object} http.ErrorResponse "Product not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /admin/products/{id} [put]
func (h *AdminHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	updProd, err := h.productSrv.UpdateProduct(r.Context(), id, req)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updProd)
}

// Orders

// ListAllOrder godoc
// @Summary Get all orders (Admin)
// @Description Returns a list of all orders with status filter and pagination
// @Tags admin
// @Accept json
// @Produce json
// @Param status query string false "Order status (pending, paid, shipped, canceled, delivered)"
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} order.OrderListResponse
// @Failure 400 {object} http.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /admin/orders [get]
func (h *AdminHandler) ListAllOrder(w http.ResponseWriter, r *http.Request) {

	filter := orderSrv.OrderFilter{
		Limit:  20,
		Offset: 0,
	}

	filter.Status = r.URL.Query().Get("status")

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(w, http.StatusBadRequest, "Invalid limit")
			return
		}
		filter.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "Invalid Offset")
			return
		}
		filter.Offset = offset
	}

	orders, err := h.orderSrv.ListOrder(r.Context(), filter)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, orders)
}

// UpdateOrderStatus godoc
// @Summary Update order status (Admin)
// @Description Changes the status of a specific order
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body order.UpdateOrderRequest true "New order status"
// @Success 200 {object} models.Order
// @Failure 400 {object} http.ErrorResponse "Invalid request body or ID"
// @Failure 404 {object} http.ErrorResponse "Order not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /admin/orders/{id}/status [put]
func (h *AdminHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}

	var req orderSrv.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updOrders, err := h.orderSrv.UpdateOrderStatus(r.Context(), id, req)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updOrders)
}

// Users

// GetUser godoc
// @Summary Get user by ID (Admin)
// @Description Returns data of a specific user
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} http.ErrorResponse "Invalid user ID"
// @Failure 404 {object} http.ErrorResponse "User not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}

	user, err := h.userSrv.GetUser(r.Context(), id)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// ListUsers godoc
// @Summary Get list of users (Admin)
// @Description Returns a list of users with filter and pagination
// @Tags admin
// @Accept json
// @Produce json
// @Param search query string false "Search query by name or email"
// @Param role query string false "User role (user, admin)"
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} user.UserListResponse
// @Failure 400 {object} http.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {

	filter := userSrv.UserFilter{
		Limit:  20,
		Offset: 0,
	}

	filter.Search = r.URL.Query().Get("search")

	if role := r.URL.Query().Get("role"); role != "" {
		filter.Role = &role
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(w, http.StatusBadRequest, "Invalid limit")
			return
		}
		filter.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "Invalid Offset")
			return
		}
		filter.Offset = offset
	}

	response, err := h.userSrv.ListUser(r.Context(), filter)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, response)
}

// DeleteUser godoc
// @Summary User deletion (Admin)
// @Description Deletes a specific user
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} http.ErrorResponse "Invalid user ID"
// @Failure 404 {object} http.ErrorResponse "User not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userSrv.DeleteUser(r.Context(), id); err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

// UpdateUser godoc
// @Summary User update (Admin)
// @Description Updates data of a specific user
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body user.UpdateUserRequest true "Data for update"
// @Success 200 {object} models.User
// @Failure 400 {object} http.ErrorResponse "Invalid request body or ID"
// @Failure 404 {object} http.ErrorResponse "User not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /admin/users/{id} [put]
func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req userSrv.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updUser, err := h.userSrv.UpdateUser(r.Context(), id, req, true)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updUser)
}
