package handler_verifications

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type VerificationDTORequest struct {
	ProblemID    string  `json:"problem_id"`
	Code         string  `json:"verification_code"`
	Language     string  `json:"language"`
	CodeTemplate *string `json:"code_template"`
	TestsNumber  int     `json:"tests_number"`
}

type VerificationDTOResponse struct {
	ID           string    `json:"id"`
	ProblemID    string    `json:"problem_id"`
	Code         string    `json:"verification_code"`
	Language     string    `json:"language"`
	CreatedAt    time.Time `json:"created_at"`
	TestsNumber  int       `json:"tests_number"`
	CodeTemplate *string   `json:"code_template"`
}

func fromModelToDTO(m model.Verification) VerificationDTOResponse {
	var codeTemplate string
	if m.CodeTemplate != nil {
		codeTemplate = m.CodeTemplate.String()
	}

	return VerificationDTOResponse{
		ID:           m.ID,
		ProblemID:    m.ProblemID,
		Code:         m.Code.String(),
		Language:     m.Language.String(),
		CreatedAt:    m.CreatedAt,
		TestsNumber:  m.TestsNumber,
		CodeTemplate: &codeTemplate,
	}
}
