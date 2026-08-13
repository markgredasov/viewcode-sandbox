package handler_verifications

import (
	"fmt"
	"net/http"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	utils_http "github.com/markgredasov/viewcode-sandbox/internal/utils/http"
)

func (h *VerificationsHTTPHandler) GetByProblemIDAndLanguage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	problemID, err := utils_http.GetUUIDPathValue(r, "problem_id")
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	languageStr := utils_http.GetStringQueryParam(r, "language")
	if languageStr == nil {
		err = fmt.Errorf("get language query param: %w", errs.ErrInvalidArgument)
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	verification, err := h.verificationsService.GetByProblemIDAndLanguage(ctx, problemID, *languageStr)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelToDTO(verification), http.StatusOK)
}
