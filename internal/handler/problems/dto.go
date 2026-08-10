package handler_problems

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type CreateProblemDTORequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProblemDTOResponse struct {
	ProblemID   string    `json:"problem_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Languages   []string  `json:"languages"`
}

func fromModelToDTO(m model.Problem) ProblemDTOResponse {
	return ProblemDTOResponse{
		ProblemID:   m.ProblemID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		Languages:   model.ConvertFrom(m.Languages),
	}
}
