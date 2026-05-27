package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/mohit838/olario-platform-backend/cmd/seed/seeders"
	"github.com/mohit838/olario-platform-backend/internal/config"
	platformpostgres "github.com/mohit838/olario-platform-backend/internal/platform/postgres"
)

// main is the seed command entrypoint.
// It runs idempotent seeders for local development and test data.
func main() {
	configPath := flag.String("config", os.Getenv("CONFIG_PATH"), "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadFile(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := platformpostgres.Open(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := seeders.Run(ctx, db); err != nil {
		log.Fatalf("run seeders: %v", err)
	}

	log.Println("seeders completed")
}
