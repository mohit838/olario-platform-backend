package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	authapp "github.com/mohit838/olario-platform-backend/internal/application/auth"
	devapp "github.com/mohit838/olario-platform-backend/internal/application/dev"
	"github.com/mohit838/olario-platform-backend/internal/config"
)

// Dependencies contains application services needed by HTTP handlers.
// The router depends on interfaces/services, not concrete database clients.
type Dependencies struct {
	DevService *devapp.Service
	Auth       *authapp.Service
}

// NewRouter builds the HTTP routing tree.
// Keeping router construction here prevents cmd/api from knowing transport
// details and gives us one place for middleware shared by every route.
func NewRouter(log *slog.Logger, cfg config.Config, deps Dependencies) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(RequestLogger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(cfg.HTTP.RequestTimeout))
	router.Use(httprate.LimitByIP(cfg.HTTP.RateLimitRequests, cfg.HTTP.RateLimitWindow))

	router.Get("/healthz", healthHandler)

	if docsEnabled(cfg) {
		router.Group(func(r chi.Router) {
			r.Use(SwaggerBasicAuth(cfg.Docs.SwaggerUsername, cfg.Docs.SwaggerPassword))
			r.Get("/swagger", swaggerUIHandler)
			r.Get("/swagger/", swaggerUIHandler)
			r.Get("/swagger/openapi.json", openAPIHandler)
		})
	}

	router.Route("/api/v1/auth", func(r chi.Router) {
		handler := NewAuthHandler(deps.Auth)
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.Refresh)
		r.Post("/logout", handler.Logout)
	})

	if cfg.Env == "local" || cfg.Env == "development" {
		router.Route("/api/v1/dev", func(r chi.Router) {
			handler := NewDevHandler(deps.DevService)
			r.Post("/full-circle", handler.FullCircle)
		})
	}

	return router
}

func docsEnabled(cfg config.Config) bool {
	return cfg.Docs.SwaggerEnabled && (cfg.Env == "local" || cfg.Env == "development")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
