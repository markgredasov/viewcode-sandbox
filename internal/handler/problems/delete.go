package handler_problems

import (
	"fmt"
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/platform/httputils"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

// Delete godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по ее ID
// @Tags Problems
// @Produce json
// @Param 	id	 path		string  true			"Идентификатор задачи" 		Format(uuid)
// @Success 204  									"Задача успешно удалена"
// @Failure 400  {object} 	response.ErrorResponse 	"Ошибка при создании модели проблемы / ID не был передан в запросе или имеет неправильный формат"
// @Failure 404  {object} 	response.ErrorResponse 	"Задача с таким ID не была найдена"
// @Failure 500  {object} 	response.ErrorResponse 	"Внутренняя ошибка сервера"
// @Router /problems/{id} [delete]
func (h *ProblemsHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	problemID, err := httputils.GetUUIDPathValue(r, "id")
	if err != nil {
		err = fmt.Errorf("get problem id: %w", err)
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	if err = h.problemsService.Delete(ctx, problemID); err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}
