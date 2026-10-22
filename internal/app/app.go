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
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
)

type App struct {
	log         *logger.Logger
	server      *server_http.HTTPServer
	diContainer *diContainer
}

func New(ctx context.Context, log *logger.Logger) (*App, error) {
	a := &App{
		diContainer: newDIContainer(log),
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
		{"pool", a.diContainer.initPool},
		{"tx manager", a.diContainer.initTxManager},
		{"sandbox adapter", a.diContainer.initAdapterSandbox},
		{"parser adapter", a.diContainer.initAdapterParser},
		{"submissions adapter", a.diContainer.initAdapterSubmissions},
		{"verifications repository", a.diContainer.initVerificationsRepository},
		{"verifications service", a.diContainer.initVerificationsService},
		{"verifications handler", a.diContainer.initVerificationsHandler},
		{"executions service", a.diContainer.initExecutionsService},
		{"executions handler", a.diContainer.initExecutionsHandler},
		{"problems repository", a.diContainer.initProblemsRepository},
		{"problems service", a.diContainer.initProblemsService},
		{"problems handler", a.diContainer.initProblemsHandler},
		{"http server", a.initHTTPServer},
	}

	for _, step := range steps {
		if err := step.fn(ctx); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

func (a *App) Run(ctx context.Context) {
	defer func() {
		closerCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()

		if err := closer.CloseAll(closerCtx); err != nil {
			a.log.Error("close deps", logger.Err(err))
		}
	}()

	if err := a.server.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("run server", logger.Err(err))
	}
}

func (a *App) initHTTPServer(ctx context.Context) error {
	apiVersion1Router := server_http.NewAPIVersionRouter(server_http.APIVersion1)
	routes, err := a.diContainer.Routes()
	if err != nil {
		return fmt.Errorf("build routes: %w", err)
	}
	apiVersion1Router.RegisterRoutes(routes...)

	log := logger.FromContext(ctx)

	serverConfig := server_http.NewConfigMust()
	server := server_http.NewHTTPServer(
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
