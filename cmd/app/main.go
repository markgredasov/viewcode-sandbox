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

// @title 		  viewcode
// @version 	  1.0
// @description   leetcode analogue by github.com/markgredasov
// @host 		  127.0.0.1:8080
// @BasePath 	  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
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
