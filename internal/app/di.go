package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/markgredasov/viewcode-sandbox/internal/app/closer"
	"github.com/markgredasov/viewcode-sandbox/internal/client/sandbox"
	sandboxconfig "github.com/markgredasov/viewcode-sandbox/internal/client/sandbox/config"
	submissionsclient "github.com/markgredasov/viewcode-sandbox/internal/client/submissions"
	executionshandler "github.com/markgredasov/viewcode-sandbox/internal/handler/executions"
	problemshandler "github.com/markgredasov/viewcode-sandbox/internal/handler/problems"
	verificationshandler "github.com/markgredasov/viewcode-sandbox/internal/handler/verifications"
	dbpostgres "github.com/markgredasov/viewcode-sandbox/internal/platform/database/postgres"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/problems"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/verifications"
	httpserver "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
	executionsservice "github.com/markgredasov/viewcode-sandbox/internal/service/executions"
	"github.com/markgredasov/viewcode-sandbox/internal/service/executions/parser"
	problemsservice "github.com/markgredasov/viewcode-sandbox/internal/service/problems"
	verificationsservice "github.com/markgredasov/viewcode-sandbox/internal/service/verifications"
)

type diContainer struct {
	log *logger.Logger

	pool      *dbpostgres.Pool
	txManager *dbpostgres.TxManager

	adapterSandbox     *sandbox.Adapter
	adapterParser      *parser.Adapter
	adapterSubmissions *submissionsclient.Adapter

	verificationsRepository *verifications.CodeVerificationRepository
	verificationsService    *verificationsservice.VerificationsService
	verificationsHandler    *verificationshandler.VerificationsHTTPHandler

	executionsService *executionsservice.ExecutionsService
	executionsHandler *executionshandler.Handler

	problemsRepository *problems.Repository
	problemsService    *problemsservice.ProblemsService
	problemsHandler    *problemshandler.ProblemsHTTPHandler
}

func newDIContainer(log *logger.Logger) *diContainer {
	return &diContainer{log: log}
}

func (d *diContainer) Pool() *dbpostgres.Pool {
	if d.pool == nil {
		panic("di: pool is not initialized — check initDeps order")
	}
	return d.pool
}

func (d *diContainer) TxManager() *dbpostgres.TxManager {
	if d.txManager == nil {
		panic("di: tx manager is not initialized — check initDeps order")
	}
	return d.txManager
}

func (d *diContainer) AdapterSandbox() *sandbox.Adapter {
	if d.adapterSandbox == nil {
		panic("di: sandbox adapter is not initialized — check initDeps order")
	}
	return d.adapterSandbox
}

func (d *diContainer) AdapterParser() *parser.Adapter {
	if d.adapterParser == nil {
		panic("di: parser adapter is not initialized — check initDeps order")
	}
	return d.adapterParser
}

func (d *diContainer) AdapterSubmissions() *submissionsclient.Adapter {
	if d.adapterSubmissions == nil {
		panic("di: submissions adapter is not initialized — check initDeps order")
	}
	return d.adapterSubmissions
}

func (d *diContainer) VerificationsRepository() *verifications.CodeVerificationRepository {
	if d.verificationsRepository == nil {
		panic("di: verifications repository is not initialized — check initDeps order")
	}
	return d.verificationsRepository
}

func (d *diContainer) ExecutionsService() *executionsservice.ExecutionsService {
	if d.executionsService == nil {
		panic("di: executions service is not initialized — check initDeps order")
	}
	return d.executionsService
}

func (d *diContainer) initPool(ctx context.Context) error {
	if d.pool != nil {
		return nil
	}

	cfg, err := dbpostgres.NewConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	pool, err := dbpostgres.NewConnectionPool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create connection pool: %w", err)
	}

	closer.Add("connection pool", func(_ context.Context) error {
		pool.Close()
		return nil
	})

	d.pool = pool
	return nil
}

func (d *diContainer) initTxManager(_ context.Context) error {
	if d.txManager != nil {
		return nil
	}
	d.txManager = dbpostgres.NewTxManager(d.Pool())
	return nil
}

func (d *diContainer) initAdapterSandbox(_ context.Context) error {
	if d.adapterSandbox != nil {
		return nil
	}

	cfg, err := sandboxconfig.NewConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	adapter, err := sandbox.NewAdapter(cfg)
	if err != nil {
		return fmt.Errorf("create sandbox adapter (host=%s, network=%s): %w",
			cfg.DockerHost, cfg.DockerNetwork, err)
	}

	d.adapterSandbox = adapter
	return nil
}

func (d *diContainer) initAdapterParser(_ context.Context) error {
	if d.adapterParser != nil {
		return nil
	}

	adapter, err := parser.NewAdapter()
	if err != nil {
		return fmt.Errorf("create parser adapter: %w", err)
	}

	d.adapterParser = adapter
	return nil
}

func (d *diContainer) initAdapterSubmissions(_ context.Context) error {
	if d.adapterSubmissions != nil {
		return nil
	}

	cfg, err := submissionsclient.NewConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	adapter, err := submissionsclient.NewAdapter(cfg,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  200 * time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   5 * time.Second,
			},
		}),
	)
	if err != nil {
		return fmt.Errorf("create submissions adapter (addr=%s): %w", cfg.Address, err)
	}

	closer.Add("submissions adapter", func(_ context.Context) error {
		return adapter.Close()
	})

	d.adapterSubmissions = adapter
	return nil
}

func (d *diContainer) initVerificationsRepository(_ context.Context) error {
	if d.verificationsRepository != nil {
		return nil
	}
	d.verificationsRepository = verifications.NewCodeVerificationRepository(d.Pool())
	return nil
}

func (d *diContainer) initVerificationsService(_ context.Context) error {
	if d.verificationsService != nil {
		return nil
	}
	d.verificationsService = verificationsservice.NewVerificationsService(d.VerificationsRepository())
	return nil
}

func (d *diContainer) initVerificationsHandler(_ context.Context) error {
	if d.verificationsHandler != nil {
		return nil
	}
	d.verificationsHandler = verificationshandler.NewVerificationsHTTPHandler(
		d.verificationsServiceOrPanic(),
	)
	return nil
}

func (d *diContainer) verificationsServiceOrPanic() *verificationsservice.VerificationsService {
	if d.verificationsService == nil {
		panic("di: verifications service is not initialized — check initDeps order")
	}
	return d.verificationsService
}

func (d *diContainer) initExecutionsService(_ context.Context) error {
	if d.executionsService != nil {
		return nil
	}
	d.executionsService = executionsservice.NewExecutionsService(
		d.AdapterSandbox(),
		d.VerificationsRepository(),
		d.AdapterParser(),
		d.AdapterSubmissions(),
	)
	return nil
}

func (d *diContainer) initExecutionsHandler(_ context.Context) error {
	if d.executionsHandler != nil {
		return nil
	}
	d.executionsHandler = executionshandler.New(d.ExecutionsService())
	return nil
}

func (d *diContainer) initProblemsRepository(_ context.Context) error {
	if d.problemsRepository != nil {
		return nil
	}
	d.problemsRepository = problems.NewProblemsRepository(d.Pool())
	return nil
}

func (d *diContainer) initProblemsService(_ context.Context) error {
	if d.problemsService != nil {
		return nil
	}
	d.problemsService = problemsservice.NewProblemsService(d.problemsRepository)
	return nil
}

func (d *diContainer) initProblemsHandler(_ context.Context) error {
	if d.problemsHandler != nil {
		return nil
	}
	d.problemsHandler = problemshandler.NewProblemsHTTPHandler(d.problemsService)
	return nil
}

func (d *diContainer) Routes() ([]httpserver.Route, error) {
	if d.executionsHandler == nil || d.problemsHandler == nil || d.verificationsHandler == nil {
		return nil, errors.New("handlers are not initialized")
	}

	var routes []httpserver.Route
	for _, getRoutes := range []func() []httpserver.Route{
		d.executionsHandler.Routes,
		d.problemsHandler.Routes,
		d.verificationsHandler.Routes,
	} {
		routes = append(routes, getRoutes()...)
	}
	return routes, nil
}
