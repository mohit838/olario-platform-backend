package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mohit838/olario-platform-backend/internal/config"
	platformpostgres "github.com/mohit838/olario-platform-backend/internal/platform/postgres"
)

// main is the migration command entrypoint.
// It runs SQL migrations against Postgres using the same config loader as the
// API, so database settings stay consistent between commands.
func main() {
	var (
		configPath     string
		migrationsPath string
	)

	fs := flag.NewFlagSet("olario-migrate", flag.ExitOnError)
	fs.StringVar(&configPath, "config", os.Getenv("CONFIG_PATH"), "path to YAML config file")
	fs.StringVar(&migrationsPath, "path", "cmd/migrate/migrations", "path to migration files")
	fs.Parse(os.Args[1:])

	if fs.NArg() != 1 {
		log.Fatal("usage: go run ./cmd/migrate --config=config/local.yml [up|down|version]")
	}

	cfg, err := config.LoadFile(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := platformpostgres.Open(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("create migrate postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer m.Close()

	switch direction := fs.Arg(0); direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("run up migrations: %v", err)
		}
		log.Println("migrations applied")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("run down migrations: %v", err)
		}
		log.Println("migrations rolled back")
	case "version":
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("read migration version: %v", err)
		}
		log.Printf("migration version=%d dirty=%t", version, dirty)
	default:
		log.Fatalf("invalid migration direction %q; use up, down, or version", direction)
	}
}
