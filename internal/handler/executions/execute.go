package handler_executions

import (
	"net/http"
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/request"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

type ExecuteDTORequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ExecuteDTOResponse struct {
	StdOut   *string        `json:"stdout"`
	StdErr   *string        `json:"stderr"`
	Duration *time.Duration `json:"duration"`
}

func (h *ExecutionsHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req ExecuteDTORequest
	if err := request.DecodeAndValidate(r, &req); err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	executionFromRequest, err := model.NewExecution(req.Language, req.Code)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	execution, err := h.executionsService.Execute(ctx, executionFromRequest)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelToDTO(execution), http.StatusOK)
}

func fromModelToDTO(m model.Execution) ExecuteDTOResponse {
	return ExecuteDTOResponse{
		StdOut:   m.StdOut,
		StdErr:   m.StdErr,
		Duration: m.Duration,
	}
}
