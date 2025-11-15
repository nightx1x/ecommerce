package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	orderSrv "github.com/nightx1x/ecommerce/interval/service/order"
	productSrv "github.com/nightx1x/ecommerce/interval/service/product"
	userSrv "github.com/nightx1x/ecommerce/interval/service/user"
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
// @Summary Створення продукту (Admin)
// @Description Створює новий продукт
// @Tags admin
// @Accept json
// @Produce json
// @Param product body productSrv.CreateProductRequest true "Дані продукту"
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
// @Summary Оновлення продукту (Admin)
// @Description Оновлює існуючий продукт
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "ID продукту"
// @Param product body productSrv.UpdateProductRequest true "Дані для оновлення"
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
// @Summary Отримати всі замовлення (Admin)
// @Description Повертає список всіх замовлень з фільтром по статусу та пагінацією
// @Tags admin
// @Accept json
// @Produce json
// @Param status query string false "Статус замовлення (pending, paid, shipped, canceled, delivered)"
// @Param limit query int false "Кількість елементів на сторінку" default(20)
// @Param offset query int false "Зміщення для пагінації" default(0)
// @Success 200 {object} orderSrv.OrderListResponse
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
// @Summary Оновлення статусу замовлення (Admin)
// @Description Змінює статус конкретного замовлення
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "ID замовлення"
// @Param status body orderSrv.UpdateOrderRequest true "Новий статус замовлення"
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
// @Summary Отримати користувача за ID (Admin)
// @Description Повертає дані конкретного користувача
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "ID користувача"
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
// @Summary Отримати список користувачів (Admin)
// @Description Повертає список користувачів з фільтром та пагінацією
// @Tags admin
// @Accept json
// @Produce json
// @Param search query string false "Пошуковий запит по імені або email"
// @Param role query string false "Роль користувача (user, admin)"
// @Param limit query int false "Кількість елементів на сторінку" default(20)
// @Param offset query int false "Зміщення для пагінації" default(0)
// @Success 200 {object} userSrv.UserListResponse
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

	response, err := h.userSrv.ListUser(r.Context(), filter) // Assuming service method is ListUsers (fix if it's ListUser)
	if err != nil {
		handlerServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Видалення користувача (Admin)
// @Description Видаляє конкретного користувача
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "ID користувача"
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
// @Summary Оновлення користувача (Admin)
// @Description Оновлює дані конкретного користувача
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "ID користувача"
// @Param user body userSrv.UpdateUserRequest true "Дані для оновлення"
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
