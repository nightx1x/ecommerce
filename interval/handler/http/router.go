package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/nightx1x/ecommerce/docs"
	authService "github.com/nightx1x/ecommerce/interval/service/auth"
	cartService "github.com/nightx1x/ecommerce/interval/service/cart"
	orderService "github.com/nightx1x/ecommerce/interval/service/order"
	productService "github.com/nightx1x/ecommerce/interval/service/product"
	userService "github.com/nightx1x/ecommerce/interval/service/user"
	httpSwagger "github.com/swaggo/http-swagger"
)

type RouterConfig struct {
	AuthService    authService.AuthService
	ProductService productService.ProductService
	CartService    cartService.CartService
	OrderService   orderService.OrderService
	UserService    userService.UserService
}

// створення путів
func NewRouter(config RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(CORS)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	r.Route("/api/v1", func(r chi.Router) {
		authHandler := NewAuthHandler(config.AuthService)
		authHandler.RegisterRoutes(r)

		ProductHandler := NewProductHandler(config.ProductService)
		ProductHandler.RegisterRoutes(r)

		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))

			userHandler := NewUserHandler(config.UserService)
			userHandler.RegisterRoutes(r)

			cartHandler := NewCartHandler(config.CartService)
			cartHandler.RegisterRoutes(r)

			orderHandler := NewOrderHandler(config.OrderService)
			orderHandler.RegisterRoutes(r)
		})

		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))
			r.Use(RequireAdmin)

			adminHandler := NewAdminHandler(
				config.ProductService,
				config.OrderService,
				config.UserService,
			)

			adminHandler.RegisterRoutes(r)
		})
	})
	return r
}
