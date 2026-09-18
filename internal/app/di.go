package app

import (
	"context"
	"os"

	"github.com/markgredasov/viewcode-sandbox/internal/app/closer"
	dbpostgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"
	handler_executions "github.com/markgredasov/viewcode-sandbox/internal/handler/executions"
	handler_problems "github.com/markgredasov/viewcode-sandbox/internal/handler/problems"
	handler_verifications "github.com/markgredasov/viewcode-sandbox/internal/handler/verifications"
	"github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/problems"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/verifications"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
	service_executions "github.com/markgredasov/viewcode-sandbox/internal/service/executions"
	"github.com/markgredasov/viewcode-sandbox/internal/service/executions/parser"
	service_problems "github.com/markgredasov/viewcode-sandbox/internal/service/problems"
	service_verifications "github.com/markgredasov/viewcode-sandbox/internal/service/verifications"
	"go.uber.org/zap"
)

type diContainer struct {
	pool *dbpostgres.ConnectionPool

	adapterSandbox *sandbox.Adapter
	adapterParser  *parser.Adapter

	verificationsRepository *verifications.CodeVerificationRepository
	verificationsService    *service_verifications.VerificationsService
	verificationsHandler    *handler_verifications.VerificationsHTTPHandler
	executionsService       *service_executions.ExecutionsService
	executionsHandler       *handler_executions.ExecutionsHTTPHandler
	problemsRepository      *problems.Repository
	problemsService         *service_problems.ProblemsService
	problemsHandler         *handler_problems.ProblemsHTTPHandler

	routes []server_http.Route
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

func (d *diContainer) AdapterSandbox() *sandbox.Adapter {
	if d.adapterSandbox == nil {
		cfg := sandbox_config.NewConfigMust()
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

func (d *diContainer) VerificationsRepository() *verifications.CodeVerificationRepository {
	if d.verificationsRepository == nil {
		d.verificationsRepository = verifications.NewCodeVerificationRepository(d.Pool())
	}

	return d.verificationsRepository
}

func (d *diContainer) VerificationsService() *service_verifications.VerificationsService {
	if d.verificationsService == nil {
		d.verificationsService = service_verifications.NewVerificationsService(d.VerificationsRepository())
	}

	return d.verificationsService
}

func (d *diContainer) VerificationsHandler() *handler_verifications.VerificationsHTTPHandler {
	if d.verificationsHandler == nil {
		d.verificationsHandler = handler_verifications.NewVerificationsHTTPHandler(d.VerificationsService())
	}

	return d.verificationsHandler
}

func (d *diContainer) ExecutionsService() *service_executions.ExecutionsService {
	if d.executionsService == nil {
		d.executionsService = service_executions.NewExecutionsService(d.AdapterSandbox(), d.VerificationsRepository(), d.AdapterParser())
	}

	return d.executionsService
}

func (d *diContainer) ExecutionsHandler() *handler_executions.ExecutionsHTTPHandler {
	if d.executionsHandler == nil {
		d.executionsHandler = handler_executions.NewExecutionsHTTPHandler(d.ExecutionsService())
	}

	return d.executionsHandler
}

func (d *diContainer) ProblemsRepository() *problems.Repository {
	if d.problemsRepository == nil {
		d.problemsRepository = problems.NewProblemsRepository(d.Pool())
	}

	return d.problemsRepository
}

func (d *diContainer) ProblemsService() *service_problems.ProblemsService {
	if d.problemsService == nil {
		d.problemsService = service_problems.NewProblemsService(d.ProblemsRepository())
	}

	return d.problemsService
}

func (d *diContainer) ProblemsHandler() *handler_problems.ProblemsHTTPHandler {
	if d.problemsHandler == nil {
		d.problemsHandler = handler_problems.NewProblemsHTTPHandler(d.ProblemsService())
	}

	return d.problemsHandler
}

func (d *diContainer) Routes() []server_http.Route {
	if len(d.routes) == 0 {
		var routes []server_http.Route

		getRoutesFuncs := []func() []server_http.Route{
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
