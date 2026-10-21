package handler_verifications

import (
	"fmt"
	"net/http"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/httputils"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

// GetByProblemIDAndLanguage godoc
// @Summary 	Найти верификацию по ProblemID и языку
// @Description Ищет верификацию по ID проблемы и выбранному языку программирования
// @Tags 		Verifications
// @Produce 	json
// @Param 		language						query		string						true	"Язык программирования (go, py, java, cpp, js, ts, csharp, ko, sw, rust)"
// @Param 		problem_id	 					path		string  					true	"Идентификатор задачи" 		Format(uuid)
// @Success		200  							{object} 	VerificationDTOResponse 			"Найденная верификация"
// @Failure 	400  							{object} 	response.ErrorResponse 				"Ошибка при валидации запроса"
// @Failure 	404  							{object} 	response.ErrorResponse 				"Верификация для выбранного языка/проблемы не существует"
// @Failure 	500  							{object} 	response.ErrorResponse 				"Внутренняя ошибка сервера"
// @Router 		/verifications/{problem_id} 	[get]
func (h *VerificationsHTTPHandler) GetByProblemIDAndLanguage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	problemID, err := httputils.GetUUIDPathValue(r, "problem_id")
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	languageStr := httputils.GetStringQueryParam(r, "language")
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
