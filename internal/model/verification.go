package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

type Verification struct {
	ID               string
	ProblemID        string
	Code             Code
	Language         Language
	CreatedAt        time.Time
	SubmissionsCount int64
}

func NewVerification(id, problemID string, codeVal string, languageVal string, createdAt time.Time, submissionsCount int64) (Verification, error) {
	if _, err := uuid.Parse(id); err != nil {
		return Verification{}, fmt.Errorf("parse id: %w", errs.ErrInvalidArgument)
	}

	if _, err := uuid.Parse(problemID); err != nil {
		return Verification{}, fmt.Errorf("parse problemID: %w", errs.ErrInvalidArgument)
	}

	code, err := NewCode(codeVal)
	if err != nil {
		return Verification{}, err
	}

	language, err := NewLanguage(languageVal)
	if err != nil {
		return Verification{}, err
	}

	return Verification{
		ID:               id,
		ProblemID:        problemID,
		Code:             code,
		Language:         language,
		CreatedAt:        createdAt,
		SubmissionsCount: submissionsCount,
	}, nil
}

func NewVerificationUnitialized(problemID, codeVal, languageVal string) (Verification, error) {
	if _, err := uuid.Parse(problemID); err != nil {
		return Verification{}, fmt.Errorf("parse problemID: %w", errs.ErrInvalidArgument)
	}

	code, err := NewCode(codeVal)
	if err != nil {
		return Verification{}, err
	}

	language, err := NewLanguage(languageVal)
	if err != nil {
		return Verification{}, err
	}

	return Verification{
		ProblemID: problemID,
		Code:      code,
		Language:  language,
	}, nil
}

func NewVerificationByCodeValue(codeVal string) (Verification, error) {
	code, err := NewCode(codeVal)
	if err != nil {
		return Verification{}, err
	}

	return Verification{
		Code: code,
	}, nil
}
