package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	orderSrv "github.com/nightx1x/ecommerce/interval/service/order"
)

type OrderHandler struct {
	OrderSrv orderSrv.OrderService
}

func NewOrderHandler(srv orderSrv.OrderService) *OrderHandler {
	return &OrderHandler{OrderSrv: srv}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/orders/{id}", h.GetOrder)
		r.Get("/orders", h.ListOrder)
		r.Post("/orders", h.CreateOrder)
		r.Put("/orders/{id}/cancel", h.CancelOrder)
	})
}

// GetOrder godoc
// @Summary Отримати замовлення за ID
// @Description Повертає детальну інформацію про замовлення за його унікальним ідентифікатором
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} models.Order
// @Failure 400 {object} http.ErrorResponse "Invalid ID"
// @Failure 404 {object} http.ErrorResponse "Order not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}

	order, err := h.OrderSrv.GetOrder(r.Context(), id)
	if err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, order)
}

// CreateOrder godoc
// @Summary Створити замовлення
// @Description Створює нове замовлення для авторизованого користувача
// @Tags orders
// @Accept json
// @Produce json
// @Param order body order.CreateOrderRequest true "Дані замовлення"
// @Success 201 {object} models.Order
// @Failure 400 {object} http.ErrorResponse "Invalid request body"
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	var req orderSrv.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := h.OrderSrv.CreateOrder(r.Context(), userID, req)
	if err != nil {
		fmt.Printf("handler lvl CreateOrder error: %+v\n", err)
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, order)
}

// CancelOrder godoc
// @Summary Скасувати замовлення
// @Description Скасовує замовлення за ID (якщо воно ще не доставлене)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} map[string]string "Order canceled message"
// @Failure 400 {object} http.ErrorResponse "Invalid ID"
// @Failure 409 {object} http.ErrorResponse "Cannot cancel order"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /orders/{id}/cancel [put]
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}

	if err := h.OrderSrv.CancelOrder(r.Context(), id); err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Order canceled!",
	})
}

// ListOrder godoc
// @Summary Отримати список замовлень
// @Description Повертає список замовлень користувача з можливістю фільтрації за статусом та пагінацією
// @Tags orders
// @Accept json
// @Produce json
// @Param status query string false "Фільтр по статусу" Enums(pending, paid, shipped, canceled, delivered)
// @Param limit query int false "Кількість елементів на сторінку" default(20)
// @Param offset query int false "Зміщення для пагінації" default(0)
// @Success 200 {object} order.OrderListResponse
// @Failure 400 {object} http.ErrorResponse "Invalid parameters"
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /orders [get]
func (h *OrderHandler) ListOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	filter := orderSrv.OrderFilter{
		UserID: &userID,
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

	orders, err := h.OrderSrv.ListOrder(r.Context(), filter)
	if err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, orders)
}

func handlerOrderError(w http.ResponseWriter, err error) {
	switch err {
	case orderSrv.ErrOrderNotFound:
		respondError(w, http.StatusNotFound, err.Error())

	case orderSrv.ErrOrderIDRequired,
		orderSrv.ErrUserIDRequired,
		orderSrv.ErrShippingAddrReq,
		orderSrv.ErrInvalidPayment,
		orderSrv.ErrProductIDRequired,
		orderSrv.ErrStatusRequired,
		orderSrv.ErrInvalidStatus,
		orderSrv.ErrOrderEmpty:
		respondError(w, http.StatusBadRequest, err.Error())

	case orderSrv.ErrOrderAlreadyCanceled,
		orderSrv.ErrCannotCancelDelivered:
		respondError(w, http.StatusConflict, err.Error())

	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
