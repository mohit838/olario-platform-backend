package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mohit838/olario-platform-backend/cmd/internal/test"
)

func AppRouters(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	testRepository := test.NewRepository(db)
	testService := test.NewService(testRepository)
	testHandler := test.NewHandler(testService)

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
	// NOTE:: Test routes for testing the architecture of the project.
	// This will be removed later.
	r.Route("/tests", func(r chi.Router) {
		r.Get("/", testHandler.GetTests)
		r.Post("/", testHandler.CreateTest)
	})

	return r
}
