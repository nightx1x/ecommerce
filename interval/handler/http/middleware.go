package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	authService "github.com/nightx1x/ecommerce/interval/service/auth"
)

type contextKey string

const (
	ContextKeyUserID    contextKey = "user_id"
	ContextKeyUserRole  contextKey = "user_role"
	ContextKeyUserEmail contextKey = "user_email"
)

func RequireAuth(authSrv authService.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authoriztion")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "Missing auth header")
				return
			}

			parts := strings.Split(authHeader, "")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, http.StatusUnauthorized, "Invalid auth header format")
				return
			}

			tokenString := parts[1]

			claims, err := authSrv.ValidateToken(tokenString)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyUserRole, claims.Role)
			ctx = context.WithValue(ctx, ContextKeyUserEmail, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(ContextKeyUserRole).(string)
		if !ok {
			respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		if role != "amin" {
			respondError(w, http.StatusForbidden, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Acces-Control-Allow-Origin", "*")
		w.Header().Set("Acces-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTION")
		w.Header().Set("Acces-Control-Allow-Headers", "Accpet,Content-Type,Content-Length,Accept-Encoding,Authorization")
		w.Header().Set("Acces-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	return userID, ok
}
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyUserRole).(string)
	return role, ok
}
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ContextKeyUserRole).(string)
	return email, ok
}
