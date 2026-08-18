package sandbox

import (
	"context"
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

func (s *sandbox) executeTests(ctx context.Context, request executionRequest) (*executionResult, error) {
	var result *executionResult
	switch request.Language {
	case "go":
		res, err := s.executeTestsLanguageGo(ctx, request)
		if err != nil {
			return result, err
		}
		result = res
	default:
		return nil, fmt.Errorf("language not provided yet: %w", errs.ErrNotImplemented)
	}

	return result, nil
}
