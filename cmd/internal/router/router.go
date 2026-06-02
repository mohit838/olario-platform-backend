package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// NOTE:: From CHI router documentations
	// Ref: https://github.com/go-chi/chi
	r.Use(middleware.RequestID)
	// Extract client IP from RemoteAddr
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Routers and router groups will be here

	return r
}
