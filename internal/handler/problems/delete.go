package handler_problems

import (
	"fmt"
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	utils_http "github.com/markgredasov/viewcode-sandbox/internal/utils/http"
)

func (h *ProblemsHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	problemID, err := utils_http.GetUUIDPathValue(r, "id")
	if err != nil {
		err = fmt.Errorf("get problem id: %w", err)
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	if err := h.problemsService.Delete(ctx, problemID); err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}
