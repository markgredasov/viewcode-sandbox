package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/markgredasov/viewcode-sandbox/internal/app"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
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

	logConfig := logger.NewConfigMust()
	log, err := logger.NewLogger(logConfig)
	if err != nil {
		fmt.Println("failed to initialize logger:", err)
		return
	}

	ctx = log.InContext(ctx)

	app.New(ctx, log).Run(ctx)
}
