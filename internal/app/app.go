package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/app/closer"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/middleware"
	httpserver "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
)

type App struct {
	server      *httpserver.HTTPServer
	diContainer *diContainer
}

func New(ctx context.Context, log *logger.Logger) (*App, error) {
	a := &App{
		diContainer: newDIContainer(ctx, log),
	}

	if err := a.diContainer.Warmup(); err != nil {
		return nil, fmt.Errorf("warmup: %w", err)
	}

	if err := a.initDeps(ctx); err != nil {
		return nil, fmt.Errorf("init deps: %w", err)
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	steps := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"http server", a.initHTTPServer},
	}

	for _, step := range steps {
		if err := step.fn(ctx); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	log := a.diContainer.Logger()

	defer func() {
		closerCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()

		if err := closer.CloseAll(closerCtx, log); err != nil {
			log.Error("close deps", logger.Err(err))
		}
	}()

	ctx = logger.InContext(ctx, log)

	if err := a.server.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	apiVersion1Router := httpserver.NewAPIVersionRouter(httpserver.APIVersion1)
	routes, err := a.diContainer.Routes()
	if err != nil {
		return fmt.Errorf("build routes: %w", err)
	}
	apiVersion1Router.RegisterRoutes(routes...)

	log := a.diContainer.Logger()

	serverConfig := httpserver.NewConfigMust()
	server := httpserver.NewHTTPServer(
		serverConfig,
		log,
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Recovery(),
		middleware.Trace(),
		middleware.CORS(),
	)

	server.RegisterAPIRouters(apiVersion1Router)
	server.RegisterSwagger()
	a.server = server

	return nil
}
