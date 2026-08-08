package app

import (
	"context"
	"fmt"
	"os"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
	"go.uber.org/zap"
)

type App struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	logger     *logger.Logger
	server     *server_http.HTTPServer
}

func NewApp(ctx context.Context, cancelFunc context.CancelFunc) *App {
	return &App{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

func (a *App) Run() error {
	if err := a.initLogger(); err != nil {
		//nolint:forbidigo // logger was not initialized to write here smth
		fmt.Println("failed to initialize logger:", err)
		os.Exit(1)
	}

	a.logger.Debug("initializing features")
	routes, err := a.initFeatures()
	if err != nil {
		a.logger.Fatal("initialize features", zap.Error(err))
	}

	a.logger.Debug("initializing HTTP server")
	if err := a.initHTTPServer(routes); err != nil {
		a.logger.Fatal("initialize HTTP server", zap.Error(err))
	}

	return a.server.Run(a.ctx)
}

func (a *App) Close() {
	a.logger.Close()
}
