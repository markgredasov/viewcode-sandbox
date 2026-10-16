package app

import (
	"context"
	"fmt"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"go.uber.org/zap"

	"github.com/markgredasov/viewcode-sandbox/internal/app/closer"
	dbpostgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"
	executionshandler "github.com/markgredasov/viewcode-sandbox/internal/handler/executions"
	problemshandler "github.com/markgredasov/viewcode-sandbox/internal/handler/problems"
	verificationshandler "github.com/markgredasov/viewcode-sandbox/internal/handler/verifications"
	"github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox"
	sandboxconfig "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	submissionsingegration "github.com/markgredasov/viewcode-sandbox/internal/integration/submissions"
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/problems"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/verifications"
	httpserver "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
	executionsservice "github.com/markgredasov/viewcode-sandbox/internal/service/executions"
	"github.com/markgredasov/viewcode-sandbox/internal/service/executions/parser"
	problemsservice "github.com/markgredasov/viewcode-sandbox/internal/service/problems"
	verificationsservice "github.com/markgredasov/viewcode-sandbox/internal/service/verifications"
)

type diContainer struct {
	pool      *dbpostgres.ConnectionPool
	txManager trm.Manager

	adapterSandbox     *sandbox.Adapter
	adapterParser      *parser.Adapter
	adapterSubmissions *submissionsingegration.Adapter

	verificationsRepository *verifications.CodeVerificationRepository
	verificationsService    *verificationsservice.VerificationsService
	verificationsHandler    *verificationshandler.VerificationsHTTPHandler
	executionsService       *executionsservice.ExecutionsService
	executionsHandler       *executionshandler.Handler
	problemsRepository      *problems.Repository
	problemsService         *problemsservice.ProblemsService
	problemsHandler         *problemshandler.ProblemsHTTPHandler

	routes []httpserver.Route
}

func newDIContainer() *diContainer {
	d := &diContainer{}

	return d
}

func (d *diContainer) initPool(ctx context.Context) {
	log := logger.FromContext(ctx)
	config := dbpostgres.NewConfigMust()
	pool, err := dbpostgres.NewConnectionPool(ctx, config)
	if err != nil {
		log.Error("initialize pool", zap.Error(err))
		os.Exit(1)
	}

	closer.Add("connection pool", func(_ context.Context) error {
		pool.Close()
		return nil
	})

	d.pool = pool
}

func (d *diContainer) Pool() *dbpostgres.ConnectionPool {
	if d.pool == nil {
		panic("connection pool must be initialized with ctx before all initializations in container")
	}

	return d.pool
}

func (d *diContainer) TxManager() trm.Manager {
	if d.txManager == nil {
		factory := trmpgx.NewDefaultFactory(d.Pool())

		m, err := manager.New(factory)
		if err != nil {
			err = fmt.Errorf("init transaction manager: %w", err)
			panic(err)
		}

		d.txManager = m
	}

	return d.txManager
}

func (d *diContainer) AdapterSandbox() *sandbox.Adapter {
	if d.adapterSandbox == nil {
		cfg := sandboxconfig.NewConfigMust()
		adapter, err := sandbox.NewAdapter(cfg)
		if err != nil {
			panic(err)
		}
		d.adapterSandbox = adapter
	}

	return d.adapterSandbox
}

func (d *diContainer) AdapterParser() *parser.Adapter {
	if d.adapterParser == nil {
		adapter, err := parser.NewAdapter()
		if err != nil {
			panic(err)
		}
		d.adapterParser = adapter
	}

	return d.adapterParser
}

func (d *diContainer) AdapterSubmissions() *submissionsingegration.Adapter {
	if d.adapterSubmissions == nil {
		cfg := submissionsingegration.NewConfigMust()
		adapter, err := submissionsingegration.NewAdapter(cfg)
		if err != nil {
			panic(err)
		}
		d.adapterSubmissions = adapter
	}

	return d.adapterSubmissions
}

func (d *diContainer) VerificationsRepository() *verifications.CodeVerificationRepository {
	if d.verificationsRepository == nil {
		d.verificationsRepository = verifications.NewCodeVerificationRepository(d.Pool(), trmpgx.DefaultCtxGetter)
	}

	return d.verificationsRepository
}

func (d *diContainer) VerificationsService() *verificationsservice.VerificationsService {
	if d.verificationsService == nil {
		d.verificationsService = verificationsservice.NewVerificationsService(d.VerificationsRepository())
	}

	return d.verificationsService
}

func (d *diContainer) VerificationsHandler() *verificationshandler.VerificationsHTTPHandler {
	if d.verificationsHandler == nil {
		d.verificationsHandler = verificationshandler.NewVerificationsHTTPHandler(d.VerificationsService())
	}

	return d.verificationsHandler
}

func (d *diContainer) ExecutionsService() *executionsservice.ExecutionsService {
	if d.executionsService == nil {
		d.executionsService = executionsservice.NewExecutionsService(d.AdapterSandbox(), d.VerificationsRepository(),
			d.AdapterParser(), d.AdapterSubmissions())
	}

	return d.executionsService
}

func (d *diContainer) ExecutionsHandler() *executionshandler.Handler {
	if d.executionsHandler == nil {
		d.executionsHandler = executionshandler.New(d.ExecutionsService())
	}

	return d.executionsHandler
}

func (d *diContainer) ProblemsRepository() *problems.Repository {
	if d.problemsRepository == nil {
		d.problemsRepository = problems.NewProblemsRepository(d.Pool())
	}

	return d.problemsRepository
}

func (d *diContainer) ProblemsService() *problemsservice.ProblemsService {
	if d.problemsService == nil {
		d.problemsService = problemsservice.NewProblemsService(d.ProblemsRepository())
	}

	return d.problemsService
}

func (d *diContainer) ProblemsHandler() *problemshandler.ProblemsHTTPHandler {
	if d.problemsHandler == nil {
		d.problemsHandler = problemshandler.NewProblemsHTTPHandler(d.ProblemsService())
	}

	return d.problemsHandler
}

func (d *diContainer) Routes() []httpserver.Route {
	if len(d.routes) == 0 {
		var routes []httpserver.Route

		getRoutesFuncs := []func() []httpserver.Route{
			d.ExecutionsHandler().Routes,
			d.ProblemsHandler().Routes,
			d.VerificationsHandler().Routes,
		}

		for _, getRoutes := range getRoutesFuncs {
			routes = append(routes, getRoutes()...)
		}

		d.routes = routes
	}

	return d.routes
}
