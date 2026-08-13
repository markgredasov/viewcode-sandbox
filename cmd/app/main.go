package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/markgredasov/viewcode-sandbox/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	app := app.NewApp(ctx, cancel)

	if err := app.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("run application:", err)
		}
	}
}
