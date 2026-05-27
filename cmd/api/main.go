package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohit838/olario-platform-backend/internal/app"
)

// main is the API process entrypoint.
// It creates the root context that is cancelled when the OS asks the process to
// stop, which lets the app shut down without dropping active HTTP requests.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "application stopped with error: %v\n", err)
		os.Exit(1)
	}
}
