package handler_problems

import (
	"context"
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
)

type ProblemsHTTPHandler struct {
	problemsService ProblemsService
}

type ProblemsService interface {
	Create(ctx context.Context, name, description string) (model.Problem, error)
}

func NewProblemsHTTPHandler(problemsService ProblemsService) *ProblemsHTTPHandler {
	return &ProblemsHTTPHandler{
		problemsService: problemsService,
	}
}

func (h *ProblemsHTTPHandler) Routes() []server_http.Route {
	return []server_http.Route{
		{
			Method:  http.MethodPost,
			Path:    "/problems",
			Handler: h.Create,
		},
	}
}
