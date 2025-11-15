package http

import (
	"encoding/json"
	"net/http"

	cartSrv "github.com/nightx1x/ecommerce/internal/service/cart"
	userSrv "github.com/nightx1x/ecommerce/internal/service/user"
)

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, ErrorResponse{
		Error:  msg,
		Status: status,
	})
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// handlerServiceUserError returns service errors
func HandlerServiceUserError(w http.ResponseWriter, err error) {
	switch err {
	case userSrv.ErrUserNotFound:
		respondError(w, http.StatusNotFound, err.Error())

	case userSrv.ErrUserIDRequired,
		userSrv.ErrInvalidEmail,
		userSrv.ErrInvalidRole,
		userSrv.ErrPasswordRequired,
		userSrv.ErrNoFields:
		respondError(w, http.StatusBadRequest, err.Error())

	case userSrv.ErrEmailAlreadyExists:
		respondError(w, http.StatusConflict, err.Error())

	default:
		respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}
func HandlerCartError(w http.ResponseWriter, err error) {
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
