package handler_verifications

import (
	"context"
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
)

type VerificationsHTTPHandler struct {
	verificationsService VerificationsService
}

type VerificationsService interface {
	Create(ctx context.Context, verification model.Verification, format *string) (model.Verification, error)
	GetByProblemIDAndLanguage(ctx context.Context, problemID, language string) (model.Verification, error)
}

func NewVerificationsHTTPHandler(verificationsService VerificationsService) *VerificationsHTTPHandler {
	return &VerificationsHTTPHandler{
		verificationsService: verificationsService,
	}
}

func (h *VerificationsHTTPHandler) Routes() []server_http.Route {
	return []server_http.Route{
		{
			Method:  http.MethodPost,
			Path:    "/verifications",
			Handler: h.Create,
		},
		{
			Method:  http.MethodGet,
			Path:    "/verifications/{problem_id}",
			Handler: h.GetByProblemIDAndLanguage,
		},
	}
}
