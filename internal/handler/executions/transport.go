package handler_executions

import (
	"context"
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/model/execution"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
)

type Handler struct {
	executionsService ExecutionsService
}

type ExecutionsService interface {
	Execute(ctx context.Context, userID, problemID, language, code string) (execution.Execution, error)
}

func New(executionsService ExecutionsService) *Handler {
	return &Handler{
		executionsService: executionsService,
	}
}

func (h *Handler) Routes() []server_http.Route {
	return []server_http.Route{
		{
			Method:  http.MethodPost,
			Path:    "/executions",
			Handler: h.Execute,
		},
	}
}
