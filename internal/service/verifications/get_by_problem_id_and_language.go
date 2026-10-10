package service_verifications

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *VerificationsService) GetByProblemIDAndLanguage(ctx context.Context, problemID, language string) (model.Verification, error) {
	languageValueObj, err := model.NewLanguage(language)
	if err != nil {
		return model.Verification{}, err
	}

	if _, err = uuid.Parse(problemID); err != nil {
		return model.Verification{}, fmt.Errorf("invalid uuid: %w: %w", err, errs.ErrInvalidArgument)
	}

	return s.verificationsRepository.GetByProblemIDAndLanguage(ctx, problemID, languageValueObj)
}
