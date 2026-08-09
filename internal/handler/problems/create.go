package handler_problems

import (
	"net/http"
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/request"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

type CreateProblemDTORequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProblemDTOResponse struct {
	ProblemID   string    `json:"problem_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

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

func fromModelToDTO(m model.Problem) CreateProblemDTOResponse {
	return CreateProblemDTOResponse{
		ProblemID:   m.ProblemID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}
