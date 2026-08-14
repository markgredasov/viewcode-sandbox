package handler_problems

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/request"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

// Create 		godoc
// @Summary 	Создать задачу
// @Description Создает задачу на viewcode
// @Tags 		Problems
// @Accept 		json
// @Produce 	json
// @Param 		RequestBody		body		CreateProblemDTORequest		true	"Request body"
// @Success		201  			{object} 	ProblemDTOResponse 					"Созданная задача"
// @Failure 	400  			{object} 	response.ErrorResponse 				"Ошибка при валидации запроса"
// @Failure 	409  			{object} 	response.ErrorResponse 				"Задача с таким названием уже существует"
// @Failure 	500  			{object} 	response.ErrorResponse 				"Внутренняя ошибка сервера"
// @Router 		/problems 		[post]
func (h *ProblemsHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req CreateProblemDTORequest
	if err := request.DecodeAndValidate(r, &req); err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	createdProblem, err := h.problemsService.Create(ctx, req.Name, req.Description)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(fromModelToDTO(createdProblem), http.StatusCreated)
}
