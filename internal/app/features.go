package app

import (
	handler_executions "github.com/markgredasov/viewcode-sandbox/internal/handler/executions"
	handler_problems "github.com/markgredasov/viewcode-sandbox/internal/handler/problems"
	"github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/problems"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/verification"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
	service_executions "github.com/markgredasov/viewcode-sandbox/internal/service/executions"
	service_problems "github.com/markgredasov/viewcode-sandbox/internal/service/problems"
)

func (a *App) initFeatures() ([]server_http.Route, error) {
	var routes []server_http.Route

	initFuncs := []func() ([]server_http.Route, error){
		a.initFeatureExecutions,
		a.initFeatureProblems,
	}

	for _, initFunc := range initFuncs {
		featureRoutes, err := initFunc()
		if err != nil {
			return routes, err
		}
		routes = append(routes, featureRoutes...)
	}

	return routes, nil

}

func (a *App) initFeatureExecutions() ([]server_http.Route, error) {
	adapterConfig := sandbox_config.NewConfigMust()
	adapter, err := sandbox.NewAdapter(adapterConfig)
	if err != nil {
		return []server_http.Route{}, err
	}

	repository := verification.NewCodeVerificationRepository(a.pool)
	service := service_executions.NewExecutionsService(adapter, repository)
	handler := handler_executions.NewExecutionsHTTPHandler(service)

	return handler.Routes(), nil
}

func (a *App) initFeatureProblems() ([]server_http.Route, error) {
	repository := problems.NewProblemsRepository(a.pool)
	service := service_problems.NewProblemsService(repository)
	handler := handler_problems.NewProblemsHTTPHandler(service)

	return handler.Routes(), nil
}
