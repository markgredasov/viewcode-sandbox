package handler_problems

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	_ "github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	utils_http "github.com/markgredasov/viewcode-sandbox/internal/utils/http"
)

// GetAll godoc
// @Summary Получить список всех задач
// @Description Отдает список всех задач на viewcode со всей информацией о них.
// @Tags Problems
// @Produce json
// @Param 	complexity		query		string					  false		"Уровень сложности задач (easy, medium, hard)"
// @Param 	categories		query		string					  false		"Категории задач (i.e. array)"
// @Success 200  			{object} 	ProblemsDTOResponse 				"Список всех задач"
// @Failure 400  			{object} 	response.ErrorResponse 				"Ошибка при создании модели проблемы"
// @Failure 500  			{object} 	response.ErrorResponse 				"Внутренняя ошибка сервера"
// @Router 	/problems 		[get]
func (h *ProblemsHTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	complexity := utils_http.GetStringQueryParam(r, "complexity")
	categories := utils_http.GetStringArrayQueryParam(r, "categories")

	problems, err := h.problemsService.GetAll(ctx, complexity, categories)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelsToDTO(problems), http.StatusOK)
}
