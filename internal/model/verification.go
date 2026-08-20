package model

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

const (
	maxCodeTemplateLength = 1_000
)

type Verification struct {
	ID               string
	ProblemID        string
	Code             Code
	Language         Language
	CreatedAt        time.Time
	SubmissionsCount int64
	CodeTemplate     *Code
}

func NewVerification(
	id,
	problemID string,
	codeVal string,
	languageVal string,
	createdAt time.Time,
	submissionsCount int64,
) (Verification, error) {
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

func (v *Verification) WithCodeTemplate(c *string) error {
	if c == nil {
		return nil
	}

	codeTemplate, err := NewCode(*c)
	if err != nil {
		return err
	}

	if utf8.RuneCountInString(codeTemplate.String()) > maxCodeTemplateLength {
		return fmt.Errorf("code template is too large: %s: %w", codeTemplate.String(), errs.ErrInvalidArgument)
	}

	v.CodeTemplate = &codeTemplate

	return nil
}
