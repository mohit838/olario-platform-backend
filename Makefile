CONFIG ?= config/config.local.yml
HTTP_ADDR ?= :8080
SWAGGER_MAIN ?= ./cmd/api/main.go
SWAGGER_OUTPUT ?= ./docs/swagger

.PHONY: help run run-local test vet tidy migrate-up migrate-down migrate-version seed docker-config docker-up docker-down dev-full-circle swagger-install swagger-generate swagger-check

help:
	@echo "Olario development commands"
	@echo "  make run                  Run API with default config"
	@echo "  make run-local            Run API with CONFIG=config/config.local.yml"
	@echo "  make test                 Run Go tests"
	@echo "  make vet                  Run go vet"
	@echo "  make tidy                 Run go mod tidy"
	@echo "  make migrate-up           Apply Postgres migrations"
	@echo "  make migrate-down         Roll back all Postgres migrations"
	@echo "  make migrate-version      Show current migration version"
	@echo "  make seed                 Seed local/test data"
	@echo "  make docker-config        Validate docker compose config"
	@echo "  make docker-up            Start docker compose services"
	@echo "  make docker-down          Stop docker compose services"
	@echo "  make dev-full-circle      Call local full-circle dev API"
	@echo "  make swagger-install      Install swag CLI for future generated Swagger"
	@echo "  make swagger-generate     Generate Swagger files after API annotation changes"
	@echo "  make swagger-check        Check Swagger starter files exist"

run:
	go run ./cmd/api --http-addr=$(HTTP_ADDR)

run-local:
	go run ./cmd/api --config=$(CONFIG) --http-addr=$(HTTP_ADDR)

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy

migrate-up:
	go run ./cmd/migrate --config=$(CONFIG) up

migrate-down:
	go run ./cmd/migrate --config=$(CONFIG) down

migrate-version:
	go run ./cmd/migrate --config=$(CONFIG) version

seed:
	go run ./cmd/seed --config=$(CONFIG)

docker-config:
	docker compose config

docker-up:
	docker compose up -d

docker-down:
	docker compose down

dev-full-circle:
	curl -i -X POST http://127.0.0.1:8080/api/v1/dev/full-circle

swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger-generate:
	swag init -g $(SWAGGER_MAIN) -o $(SWAGGER_OUTPUT) --parseInternal

swagger-check:
	test -f internal/http/swagger.go
	test -f config/config.example.yml
