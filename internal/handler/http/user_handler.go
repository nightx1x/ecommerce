package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	userSrv "github.com/nightx1x/ecommerce/internal/service/user"
)

type UserHandler struct {
	UserSrv userSrv.UserService
}

func NewUserHandler(srv userSrv.UserService) *UserHandler {
	return &UserHandler{UserSrv: srv}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/users", h.ListUser)
		r.Put("/users", h.UpdateUser)
	})

}

// ListUser godoc
// @Summary Get user information
// @Description Returns data of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 404 {object} http.ErrorResponse "User not found"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	user, err := h.UserSrv.GetUser(r.Context(), userID)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user data
// @Description Updates user data (email, password, first name, last name, role)
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.UpdateOwnUserRequest true "Data for updating user"
// @Success 200 {object} models.User
// @Failure 400 {object} http.ErrorResponse "Invalid request body or invalid fields"
// @Failure 401 {object} http.ErrorResponse "User not authorized"
// @Failure 404 {object} http.ErrorResponse "User not found"
// @Failure 409 {object} http.ErrorResponse "Email already exists"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	var req userSrv.UpdateOwnUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updReq := userSrv.UpdateUserRequest{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Role:      nil,
	}
	updUser, err := h.UserSrv.UpdateUser(r.Context(), userID, updReq, false)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updUser)
}

// handlerServiceUserError returns service errors
func handlerServiceUserError(w http.ResponseWriter, err error) {
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
