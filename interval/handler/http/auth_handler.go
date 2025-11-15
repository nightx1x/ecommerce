package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	authSrv "github.com/nightx1x/ecommerce/interval/service/auth"
)

type AuthHandler struct {
	AuthSrv authSrv.AuthService
}

func NewAuthHandler(srv authSrv.AuthService) *AuthHandler {
	return &AuthHandler{AuthSrv: srv}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)
		r.Post("/auth/refresh", h.RefreshToken)
	})
}

// Register godoc
// @Summary Реєстрація нового користувача
// @Description Реєструє нового користувача та повертає токени доступу
// @Tags auth
// @Accept json
// @Produce json
// @Param user body auth.RegisterRequest true "Дані для реєстрації"
// @Success 201 {object} auth.AuthResponse
// @Failure 400 {object} http.ErrorResponse "Invalid request body"
// @Failure 409 {object} http.ErrorResponse "User already exists"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authSrv.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.AuthSrv.Register(r.Context(), req)
	if err != nil {
		handlerAuthError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, resp)
}

// Login godoc
// @Summary Логін користувача
// @Description Логін користувача та повертає токени доступу
// @Tags auth
// @Accept json
// @Produce json
// @Param user body auth.LoginRequset true "Дані для логіну"
// @Success 200 {object} auth.AuthResponse
// @Failure 400 {object} http.ErrorResponse "Invalid request body"
// @Failure 401 {object} http.ErrorResponse "Invalid credentials"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authSrv.LoginRequset
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	resp, err := h.AuthSrv.Login(r.Context(), req)
	if err != nil {
		handlerAuthError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary Оновлення токена доступу
// @Description Оновлює access token за допомогою refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body auth.RefreshRequest true "Refresh Token"
// @Success 200 {object} auth.AuthResponse
// @Failure 400 {object} http.ErrorResponse "Invalid request body"
// @Failure 401 {object} http.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} http.ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req authSrv.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.AuthSrv.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		handlerAuthError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func handlerAuthError(w http.ResponseWriter, err error) {
	switch err {
	case authSrv.ErrInvalidCredentials:
		respondError(w, http.StatusUnauthorized, err.Error())
	case authSrv.ErrUserAlreadyExists:
		respondError(w, http.StatusConflict, err.Error())
	case authSrv.ErrInvalidToken, authSrv.ErrExpiredToken:
		respondError(w, http.StatusUnauthorized, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
