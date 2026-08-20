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
}

type VerificationDTOResponse struct {
	ID               string    `json:"id"`
	ProblemID        string    `json:"problem_id"`
	Code             string    `json:"verification_code"`
	Language         string    `json:"language"`
	CreatedAt        time.Time `json:"created_at"`
	SubmissionsCount int64     `json:"submissions_count"`
	CodeTemplate     *string   `json:"code_template"`
}

func fromModelToDTO(m model.Verification) VerificationDTOResponse {
	codeTemplate := m.CodeTemplate.String()
	return VerificationDTOResponse{
		ID:               m.ID,
		ProblemID:        m.ProblemID,
		Code:             m.Code.String(),
		Language:         m.Language.String(),
		CreatedAt:        m.CreatedAt,
		SubmissionsCount: m.SubmissionsCount,
		CodeTemplate:     &codeTemplate,
	}
}
