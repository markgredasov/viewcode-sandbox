package handler_verifications

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/request"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	utils_http "github.com/markgredasov/viewcode-sandbox/internal/utils/http"
)

// Create 		godoc
// @Summary 	Создать верификацию
// @Description Создает верификацию (добавляет тесты для проблемы для выбранного языка). В теле запроса код должен быть закодирован в стандарте base64
// @Tags 		Verifications
// @Accept 		json
// @Produce 	json
// @Param 		RequestBody		body		VerificationDTORequest		true	"Request body"
// @Param 		format			query		string						true	"Формат кодирования верификационного кода (доступен base64)"
// @Success		201  			{object} 	VerificationDTOResponse 			"Созданная верификация"
// @Failure 	400  			{object} 	response.ErrorResponse 				"Ошибка при валидации запроса"
// @Failure 	404  			{object} 	response.ErrorResponse 				"Проблема с таким ID не найдена"
// @Failure 	409  			{object} 	response.ErrorResponse 				"Верификация для выбранной задачи и языка уже существует"
// @Failure 	500  			{object} 	response.ErrorResponse 				"Внутренняя ошибка сервера"
// @Router 		/verifications 	[post]
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

	createdVerification, err := h.verificationsService.Create(ctx, verificationFromRequest, format, req.CodeTemplate)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelToDTO(createdVerification), http.StatusCreated)
}
