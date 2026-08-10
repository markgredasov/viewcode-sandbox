package handler_problems

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

func (h *ProblemsHTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	problems, err := h.problemsService.GetAll(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	response := make([]model.Problem, len(problems))

	for i := range problems {
		response[i] = model.Problem(fromModelToDTO(problems[i]))
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
