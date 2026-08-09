package app

import (
	handler_executions "github.com/markgredasov/viewcode-sandbox/internal/handler/executions"
	"github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/repository/verification"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
	service_executions "github.com/markgredasov/viewcode-sandbox/internal/service/executions"
)

func (a *App) initFeatures() ([]server_http.Route, error) {
	routes := make([]server_http.Route, 0)

	executionsRoutes, err := a.initFeatureExecutions()
	if err != nil {
		return routes, err
	}

	return append(routes, executionsRoutes...), nil
}

func (a *App) initFeatureExecutions() ([]server_http.Route, error) {
	adapterConfig := sandbox_config.NewConfigMust()
	adapter, err := sandbox.NewAdapter(adapterConfig)
	if err != nil {
		return []server_http.Route{}, err
	}

	repository := verification.NewVerificationRepository(a.pool)
	service := service_executions.NewExecutionsService(adapter, repository)
	handler := handler_executions.NewExecutionsHTTPHandler(service)

	return handler.Routes(), nil
}
