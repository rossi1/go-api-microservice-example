package handler

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API-example
// @version 1.0
// @description This is a go-api-microservice-example.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.basic- BasicAuth

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func (h *Handler) GetRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.BasicAuth("my realm", map[string]string{os.Getenv("username"): os.Getenv("password")}))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"), //The url pointing to API definition"
	))
	r.Route("/category", func(r chi.Router) {
		r.Post("/", h.createCategoryAPIView)
		r.Get("/", h.getAllCategoriesAPIView)
		r.Route("/{categoryId}", func(r chi.Router) {
			r.Get("/", h.getCategoryAPIView)
			r.Delete("/", h.deleteCategoryAPIView)
			r.Put("/", h.updateCategoryAPIView)
			r.Route("/products", func(r chi.Router) {
				r.Post("/", h.createProductAPIView)
				r.Get("/", h.getAllProductsAPIView)
				r.Route("/{productId}", func(r chi.Router) {
					r.Delete("/", h.deleteProductAPIView)
					r.Put("/", h.updateProductAPIView)
					r.Get("/", h.getProductAPIView)
				})
			})

		})

	})

	return r
}
