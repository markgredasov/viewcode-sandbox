package handler_problems

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	_ "github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

// GetAll godoc
// @Summary Получить список всех задач
// @Description Отдает список всех задач на viewcode со всей информацией о них.
// @Tags Problems
// @Produce json
// @Success 200  {object} 	ProblemsDTOResponse 	"Список всех задач"
// @Failure 400  {object} 	response.ErrorResponse 	"Ошибка при создании модели проблемы"
// @Failure 500  {object} 	response.ErrorResponse 	"Внутренняя ошибка сервера"
// @Router /problems [get]
func (h *ProblemsHTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	problems, err := h.problemsService.GetAll(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelsToDTO(problems), http.StatusOK)
}
