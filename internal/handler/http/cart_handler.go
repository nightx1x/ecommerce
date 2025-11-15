package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	cartSrv "github.com/nightx1x/ecommerce/internal/service/cart"
)

type CartHandler struct {
	CartSrv cartSrv.CartService
}

func NewCartHandler(srv cartSrv.CartService) *CartHandler {
	return &CartHandler{CartSrv: srv}
}

func (h *CartHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/cart", h.ListItems)
		r.Post("/cart/items", h.AddItem)
		r.Put("/cart/items/{id}", h.UpdateItem)
		r.Delete("/cart/items/{id}", h.DeleteItem)
		r.Delete("/cart", h.ClearCart)
	})
}

// AddItem godoc
// @Summary Adds an item to the cart
// @Description Adds a new item to the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param item body cart.AddCartItemRequest true "Item to add"
// @Success 201 {object} models.Cart
// @Failure 400 {object} http.ErrorResponse "Invalid request body or quantity"
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart/items [post]
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}
	var req cartSrv.AddCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.CartSrv.AddItem(r.Context(), userID, req)
	if err != nil {
		handlerCartError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

// UpdateItem godoc
// @Summary Updates an item in the cart
// @Description Changes the quantity or details of an item in the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart item ID"
// @Param item body cart.UpdateCartItemRequest true "Updated item data"
// @Success 200 {object} models.Cart
// @Failure 400 {object} http.ErrorResponse "Invalid request body or quantity"
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 404 {object} http.ErrorResponse "Item not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart/items/{id} [put]
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req cartSrv.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.CartSrv.UpdateItem(r.Context(), userID, id, req)
	if err != nil {
		handlerCartError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

// DeleteItem godoc
// @Summary Deletes an item from the cart
// @Description Removes a specific item from the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart item ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} http.ErrorResponse "Invalid item ID"
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 404 {object} http.ErrorResponse "Item not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart/items/{id} [delete]
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	if err := h.CartSrv.DeleteItem(r.Context(), userID, id); err != nil {
		handlerCartError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "item deleted",
	})
}

// ClearCart godoc
// @Summary Clears the cart
// @Description Removes all items from the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart [delete]
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	if err := h.CartSrv.ClearItem(r.Context(), userID); err != nil {
		handlerCartError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "cart cleared successfully",
	})
}

// ListItems godoc
// @Summary Returns the list of items in the cart
// @Description Returns all user items in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} cart.CartListResponse
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart [get]
func (h *CartHandler) ListItems(w http.ResponseWriter, r *http.Request) {

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	items, err := h.CartSrv.ListItem(r.Context(), userID)
	if err != nil {
		handlerCartError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func handlerCartError(w http.ResponseWriter, err error) {
	switch err {
	case cartSrv.ErrItemNotFound:
		respondError(w, http.StatusNotFound, err.Error())
	case cartSrv.ErrInvalidQuantity,
		cartSrv.ErrProductNotAvailable,
		cartSrv.ErrInvalidProductID:
		respondError(w, http.StatusBadRequest, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
