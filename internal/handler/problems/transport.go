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
	Create(ctx context.Context, name, description, complexity string, categories []string) (model.Problem, error)
	GetAll(ctx context.Context, complexity *string, categories *[]string) ([]model.Problem, error)
	GetByID(ctx context.Context, id string) (model.Problem, error)
	Delete(ctx context.Context, id string) error
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
		{
			Method:  http.MethodGet,
			Path:    "/problems",
			Handler: h.GetAll,
		},
		{
			Method:  http.MethodGet,
			Path:    "/problems/{id}",
			Handler: h.GetByID,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/problems/{id}",
			Handler: h.Delete,
		},
	}
}
