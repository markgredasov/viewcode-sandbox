package handler_executions

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/request"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

// Execute 		godoc
// @Summary 	Запустить код
// @Description Запускает код для выбранной задачи и языка
// @Tags 		Executions
// @Accept 		json
// @Produce 	json
// @Param 		RequestBody		body		ExecuteDTORequest			true	"Request body"
// @Success		200  			{object} 	ExecuteDTOResponse		 			"Результаты прохождения тестов"
// @Failure 	400  			{object} 	response.ErrorResponse 				"Ошибка при валидации запроса / создании модели исполнения"
// @Failure 	404  			{object} 	response.ErrorResponse 				"Верификация для выбранного языка/проблемы не существует"
// @Failure 	500  			{object} 	response.ErrorResponse 				"Внутренняя ошибка сервера"
// @Router 		/executions 	[post]
func (h *Handler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req ExecuteDTORequest
	if err := request.DecodeAndValidate(r, &req); err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	execution, err := h.executionsService.Execute(ctx, model.NewUUID().String(), req.ProblemID, req.Language, req.Code)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelToDTO(execution), http.StatusOK)
}
