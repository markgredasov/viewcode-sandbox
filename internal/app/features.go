package app

import (
	handler_executions "github.com/markgredasov/viewcode-sandbox/internal/handler/executions"
	handler_problems "github.com/markgredasov/viewcode-sandbox/internal/handler/problems"
	handler_verifications "github.com/markgredasov/viewcode-sandbox/internal/handler/verifications"
	"github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/problems"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/verifications"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
	service_executions "github.com/markgredasov/viewcode-sandbox/internal/service/executions"
	service_problems "github.com/markgredasov/viewcode-sandbox/internal/service/problems"
	service_verifications "github.com/markgredasov/viewcode-sandbox/internal/service/verifications"
)

func (a *App) initFeatures() ([]server_http.Route, error) {
	var routes []server_http.Route

	initFuncs := []func() ([]server_http.Route, error){
		a.initFeatureExecutions,
		a.initFeatureProblems,
		a.initFeatureVerifications,
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

	repository := verifications.NewCodeVerificationRepository(a.pool)
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

func (a *App) initFeatureVerifications() ([]server_http.Route, error) {
	repository := verifications.NewCodeVerificationRepository(a.pool)
	service := service_verifications.NewVerificationsService(repository)
	handler := handler_verifications.NewVerificationsHTTPHandler(service)

	return handler.Routes(), nil
}
