package handler_executions

import (
	"context"
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
)

type ExecutionsHTTPHandler struct {
	executionsService ExecutionsService
}

type ExecutionsService interface {
	Execute(ctx context.Context, exec model.Execution) (model.Execution, error)
}

func NewExecutionsHTTPHandler(executionsService ExecutionsService) ExecutionsHTTPHandler {
	return ExecutionsHTTPHandler{
		executionsService: executionsService,
	}
}

func (h *ExecutionsHTTPHandler) Routes() []server_http.Route {
	return []server_http.Route{
		{
			Method:  http.MethodPost,
			Path:    "/executions",
			Handler: h.Execute,
		},
	}
}
