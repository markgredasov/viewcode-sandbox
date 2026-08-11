package handler_verifications

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/request"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	utils_http "github.com/markgredasov/viewcode-sandbox/internal/utils/http"
)

func (h *VerificationsHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req VerificationDTORequest
	if err := request.DecodeAndValidate(r, &req); err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	verificationFromRequest, err := model.NewVerificationUnitialized(req.ProblemID, req.Code, req.Language)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	format := utils_http.GetStringQueryParam(r, "format")

	createdVerification, err := h.verificationsService.Create(ctx, verificationFromRequest, format)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelToDTO(createdVerification), http.StatusCreated)
}
