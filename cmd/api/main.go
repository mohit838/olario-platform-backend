package main

import (
	"log"
	"net/http"

	"github.com/mohit838/olario-platform-backend/cmd/internal/config"
	"github.com/mohit838/olario-platform-backend/cmd/internal/database"
	"github.com/mohit838/olario-platform-backend/cmd/internal/router"
)

func main() {
	// Load environment variables
	cfg, err := config.LoadConfig("./config/.env")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// DB connections
	db, err := database.LoadPostgresDB(cfg.AppDB)
	if err != nil {
		log.Fatalf("DB connections needed!")
	}
	defer db.Close()
	log.Println("DB is connected")

	// Add Chi router
	handler := router.NewRouter(db)

	// Start server
	log.Printf("Server start on port: %s", cfg.AppPort)
	port := ":" + cfg.AppPort
	err = http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatalf("Server is not running %v", err)
	}
}
