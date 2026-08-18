package service_executions

import (
	"context"
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *ExecutionsService) Execute(ctx context.Context, exec model.Execution) (model.Execution, error) {
	problemID := exec.ProblemID.String()

	tests, err := s.testsRepo.GetByProblemIDAndLanguage(ctx, problemID, exec.Language)
	if err != nil {
		return model.Execution{}, err
	}

	var execution model.Execution

	switch exec.Language {
	case model.LanguageGo:
		//nolint:govet // not needed
		executionGo, err := s.executeLanguageGo(ctx, exec, tests)
		if err != nil {
			return model.Execution{}, err
		}
		execution = executionGo
	case model.LanguagePython:
		fallthrough
	case model.LanguageJava:
		fallthrough
	case model.LanguageCpp:
		fallthrough
	case model.LanguageJavaScript:
		fallthrough
	case model.LanguageTypeScript:
		fallthrough
	case model.LanguageCSharp:
		fallthrough
	case model.LanguageKotlin:
		fallthrough
	case model.LanguageSwift:
		fallthrough
	case model.LanguageRust:
		fallthrough
	default:
		return model.Execution{}, fmt.Errorf("language not provided yet: %w", errs.ErrNotImplemented)
	}

	return execution, nil
}
