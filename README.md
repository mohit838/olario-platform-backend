# Olario Platform Backend

A minimal Go backend service for the Olario platform.

## Requirements

- Go `1.26.2` or newer
- Air (optional, for live reload)

## Run the app

```bash
go run ./cmd
```

## Run with Air (live reload)

Install Air (if needed):

```bash
go install github.com/air-verse/air@latest
```

Start dev server:

```bash
air
```

The Air config is in `.air.toml` and builds from `./cmd`.

## Project structure

```text
.
├── cmd/
│   └── main.go
├── go.mod
└── .air.toml
```
