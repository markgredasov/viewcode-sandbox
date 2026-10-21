package handler_problems

import (
	"fmt"
	"net/http"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/httputils"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

// GetByID godoc
// @Summary Получить список задачи по ее UUID
// @Description Отдает задачу на viewcode по ее UUID
// @Tags Problems
// @Produce json
// @Param 	id	 path		string  true			"Идентификатор задачи" 		Format(uuid)
// @Success 200  {object} 	ProblemDTOResponse 		"Задача с указанным ID"
// @Failure 400  {object} 	response.ErrorResponse 	"Ошибка при создании модели проблемы / ID не был передан в запросе или имеет неправильный формат"
// @Failure 404  {object} 	response.ErrorResponse 	"Задача с таким ID не была найдена"
// @Failure 500  {object} 	response.ErrorResponse 	"Внутренняя ошибка сервера"
// @Router /problems/{id} [get]
func (h *ProblemsHTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	problemID, err := httputils.GetUUIDPathValue(r, "id")
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
