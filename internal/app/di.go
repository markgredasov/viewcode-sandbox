package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/markgredasov/viewcode-sandbox/internal/app/closer"
	"github.com/markgredasov/viewcode-sandbox/internal/client/sandbox"
	sandboxcfg "github.com/markgredasov/viewcode-sandbox/internal/client/sandbox/config"
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

	poolOnce      func() (*dbpostgres.Pool, error)
	txManagerOnce func() (*dbpostgres.TxManager, error)

	sandboxAdapterOnce     func() (*sandbox.Adapter, error)
	parserAdapterOnce      func() (*parser.Adapter, error)
	submissionsAdapterOnce func() (*submissionsclient.Adapter, error)

	verificationsRepoOnce    func() (*verifications.CodeVerificationRepository, error)
	verificationsServiceOnce func() (*verificationsservice.VerificationsService, error)
	verificationsHandlerOnce func() (*verificationshandler.VerificationsHTTPHandler, error)
	executionsServiceOnce    func() (*executionsservice.ExecutionsService, error)
	executionsHandlerOnce    func() (*executionshandler.Handler, error)
	problemsRepoOnce         func() (*problems.Repository, error)
	problemsServiceOnce      func() (*problemsservice.ProblemsService, error)
	problemsHandlerOnce      func() (*problemshandler.ProblemsHTTPHandler, error)
}

func newDIContainer(ctx context.Context, log *logger.Logger) *diContainer {
	d := &diContainer{log: log}

	d.poolOnce = sync.OnceValues(func() (*dbpostgres.Pool, error) { return d.initPool(ctx) })
	d.txManagerOnce = sync.OnceValues(d.initTxManager)
	d.sandboxAdapterOnce = sync.OnceValues(d.initSandboxAdapter)
	d.parserAdapterOnce = sync.OnceValues(d.initParserAdapter)
	d.submissionsAdapterOnce = sync.OnceValues(d.initSubmissionsAdapter)
	d.verificationsRepoOnce = sync.OnceValues(d.initVerificationsRepo)
	d.verificationsServiceOnce = sync.OnceValues(d.initVerificationsService)
	d.verificationsHandlerOnce = sync.OnceValues(d.initVerificationsHandler)
	d.executionsServiceOnce = sync.OnceValues(d.initExecutionsService)
	d.executionsHandlerOnce = sync.OnceValues(d.initExecutionsHandler)
	d.problemsRepoOnce = sync.OnceValues(d.initProblemsRepo)
	d.problemsServiceOnce = sync.OnceValues(d.initProblemsService)
	d.problemsHandlerOnce = sync.OnceValues(d.initProblemsHandler)

	return d
}

func (d *diContainer) Logger() *logger.Logger {
	return d.log
}

func (d *diContainer) Pool() (*dbpostgres.Pool, error) {
	return d.poolOnce()
}

func (d *diContainer) TxManager() (*dbpostgres.TxManager, error) {
	return d.txManagerOnce()
}

func (d *diContainer) SandboxAdapter() (*sandbox.Adapter, error) {
	return d.sandboxAdapterOnce()
}

func (d *diContainer) ParserAdapter() (*parser.Adapter, error) {
	return d.parserAdapterOnce()
}

func (d *diContainer) SubmissionsAdapter() (*submissionsclient.Adapter, error) {
	return d.submissionsAdapterOnce()
}

func (d *diContainer) VerificationsRepo() (*verifications.CodeVerificationRepository, error) {
	return d.verificationsRepoOnce()
}

func (d *diContainer) VerificationsService() (*verificationsservice.VerificationsService, error) {
	return d.verificationsServiceOnce()
}

func (d *diContainer) VerificationsHandler() (*verificationshandler.VerificationsHTTPHandler, error) {
	return d.verificationsHandlerOnce()
}

func (d *diContainer) ExecutionsService() (*executionsservice.ExecutionsService, error) {
	return d.executionsServiceOnce()
}

func (d *diContainer) ExecutionsHandler() (*executionshandler.Handler, error) {
	return d.executionsHandlerOnce()
}

func (d *diContainer) ProblemsRepo() (*problems.Repository, error) {
	return d.problemsRepoOnce()
}

func (d *diContainer) ProblemsService() (*problemsservice.ProblemsService, error) {
	return d.problemsServiceOnce()
}

func (d *diContainer) ProblemsHandler() (*problemshandler.ProblemsHTTPHandler, error) {
	return d.problemsHandlerOnce()
}

func (d *diContainer) Routes() ([]httpserver.Route, error) {
	if err := d.Warmup(); err != nil {
		return nil, err
	}

	executionsHandler, err := d.ExecutionsHandler()
	if err != nil {
		return nil, err
	}

	problemsHandler, err := d.ProblemsHandler()
	if err != nil {
		return nil, err
	}

	verificationsHandler, err := d.VerificationsHandler()
	if err != nil {
		return nil, err
	}

	var routes []httpserver.Route
	for _, getRoutes := range []func() []httpserver.Route{
		executionsHandler.Routes,
		problemsHandler.Routes,
		verificationsHandler.Routes,
	} {
		routes = append(routes, getRoutes()...)
	}
	return routes, nil
}

func (d *diContainer) initPool(ctx context.Context) (*dbpostgres.Pool, error) {
	log := d.Logger()

	cfg := dbpostgres.NewConfigMust()
	pool, err := dbpostgres.NewConnectionPool(ctx, cfg)
	if err != nil {
		log.Error("initialize pool", logger.Err(err))
		return nil, fmt.Errorf("init pool: %w", err)
	}

	closer.Add("connection pool", func(_ context.Context) error {
		pool.Close()
		return nil
	})

	return pool, nil
}

func (d *diContainer) initTxManager() (*dbpostgres.TxManager, error) {
	pool, err := d.Pool()
	if err != nil {
		return nil, fmt.Errorf("init tx manager: %w", err)
	}

	return dbpostgres.NewTxManager(pool), nil
}

func (d *diContainer) initSandboxAdapter() (*sandbox.Adapter, error) {
	cfg, err := sandboxcfg.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("init sandbox adapter config: %w", err)
	}

	adapter, err := sandbox.NewAdapter(cfg)
	if err != nil {
		return nil, fmt.Errorf("init sandbox adapter: %w", err)
	}

	return adapter, nil
}

func (d *diContainer) initParserAdapter() (*parser.Adapter, error) {
	adapter, err := parser.NewAdapter()
	if err != nil {
		return nil, fmt.Errorf("init parser adapter: %w", err)
	}

	return adapter, nil
}

func (d *diContainer) initSubmissionsAdapter() (*submissionsclient.Adapter, error) {
	cfg, err := submissionsclient.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("init submissions adapter config: %w", err)
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
		return nil, fmt.Errorf("create submissions adapter (addr=%s): %w", cfg.Address, err)
	}

	closer.Add("submissions adapter", func(_ context.Context) error {
		return adapter.Close()
	})

	return adapter, nil
}

func (d *diContainer) initVerificationsRepo() (*verifications.CodeVerificationRepository, error) {
	pool, err := d.Pool()
	if err != nil {
		return nil, fmt.Errorf("init verifications repo: %w", err)
	}
	return verifications.NewCodeVerificationRepository(pool), nil
}

func (d *diContainer) initVerificationsService() (*verificationsservice.VerificationsService, error) {
	repo, err := d.VerificationsRepo()
	if err != nil {
		return nil, fmt.Errorf("init verifications service: %w", err)
	}

	return verificationsservice.NewVerificationsService(repo), nil
}

func (d *diContainer) initVerificationsHandler() (*verificationshandler.VerificationsHTTPHandler, error) {
	service, err := d.VerificationsService()
	if err != nil {
		return nil, fmt.Errorf("init verifications handler: %w", err)
	}
	return verificationshandler.NewVerificationsHTTPHandler(service), nil
}

func (d *diContainer) initExecutionsService() (*executionsservice.ExecutionsService, error) {
	sandboxAdapter, err := d.SandboxAdapter()
	if err != nil {
		return nil, fmt.Errorf("init executions service: sandbox adapter: %w", err)
	}

	verificationsRepo, err := d.VerificationsRepo()
	if err != nil {
		return nil, fmt.Errorf("init executions service: verifications repo: %w", err)
	}

	parserAdapter, err := d.ParserAdapter()
	if err != nil {
		return nil, fmt.Errorf("init executions service: parser adapter: %w", err)
	}

	submissionsAdapter, err := d.SubmissionsAdapter()
	if err != nil {
		return nil, fmt.Errorf("init executions service: submissions adapter: %w", err)
	}

	return executionsservice.NewExecutionsService(
		sandboxAdapter,
		verificationsRepo,
		parserAdapter,
		submissionsAdapter,
	), nil
}

func (d *diContainer) initExecutionsHandler() (*executionshandler.Handler, error) {
	service, err := d.ExecutionsService()
	if err != nil {
		return nil, fmt.Errorf("init executions handler: %w", err)
	}

	return executionshandler.New(service), nil
}

func (d *diContainer) initProblemsRepo() (*problems.Repository, error) {
	pool, err := d.Pool()
	if err != nil {
		return nil, fmt.Errorf("init problems repo: %w", err)
	}

	return problems.NewProblemsRepository(pool), nil
}

func (d *diContainer) initProblemsService() (*problemsservice.ProblemsService, error) {
	repo, err := d.ProblemsRepo()
	if err != nil {
		return nil, fmt.Errorf("init problems service: %w", err)
	}

	return problemsservice.NewProblemsService(repo), nil
}

func (d *diContainer) initProblemsHandler() (*problemshandler.ProblemsHTTPHandler, error) {
	service, err := d.ProblemsService()
	if err != nil {
		return nil, fmt.Errorf("init problems handler: %w", err)
	}

	return problemshandler.NewProblemsHTTPHandler(service), nil
}

func (d *diContainer) Warmup() error {
	var errs []error

	if _, err := d.TxManager(); err != nil {
		errs = append(errs, fmt.Errorf("tx manager: %w", err))
	}
	if _, err := d.SandboxAdapter(); err != nil {
		errs = append(errs, fmt.Errorf("sandbox adapter: %w", err))
	}
	if _, err := d.ParserAdapter(); err != nil {
		errs = append(errs, fmt.Errorf("parser adapter: %w", err))
	}
	if _, err := d.SubmissionsAdapter(); err != nil {
		errs = append(errs, fmt.Errorf("submissions adapter: %w", err))
	}
	if _, err := d.VerificationsHandler(); err != nil {
		errs = append(errs, fmt.Errorf("verifications handler: %w", err))
	}
	if _, err := d.ExecutionsHandler(); err != nil {
		errs = append(errs, fmt.Errorf("executions handler: %w", err))
	}
	if _, err := d.ProblemsHandler(); err != nil {
		errs = append(errs, fmt.Errorf("problems handler: %w", err))
	}

	return errors.Join(errs...)
}
