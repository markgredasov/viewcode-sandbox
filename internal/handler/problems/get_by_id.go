package handler_problems

import (
	"fmt"
	"net/http"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	utils_http "github.com/markgredasov/viewcode-sandbox/internal/utils/http"
)

func (h *ProblemsHTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	problemID, err := utils_http.GetUUIDPathValue(r, "id")
	if err != nil {
		err = fmt.Errorf("get problem id: %w", errs.ErrInvalidArgument)
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	problemGot, err := h.problemsService.GetByID(ctx, problemID)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelToDTO(problemGot), http.StatusOK)
}
