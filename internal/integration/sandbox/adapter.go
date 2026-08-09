package sandbox

import (
	"context"
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

var languageToExecutionLanguage = map[model.Language]string{
	model.LanguageGo:     "go",
	model.LanguagePython: "python3",
}

type Adapter struct {
	sandbox *sandbox
}

func NewAdapter(config sandbox_config.Config) (*Adapter, error) {
	sandbox, err := newSandbox(config)
	if err != nil {
		return nil, fmt.Errorf("new adapter: %w", err)
	}

	return &Adapter{
		sandbox: sandbox,
	}, nil
}

func (a *Adapter) RunTests(ctx context.Context, exec model.Execution, tests model.Verification) (model.Execution, error) {
	language, ok := languageToExecutionLanguage[exec.Language]
	if !ok {
		return model.Execution{}, errs.ErrInvalidLanguage
	}

	testsCode := tests.Code.String()

	req := executionRequest{
		Language: language,
		Code:     exec.Code.String(),
		Timeout:  a.sandbox.config.ContainerExecTimeout,
		Tests:    &testsCode,
	}

	execResult, err := a.sandbox.executeTests(ctx, req)
	if err != nil {
		return model.Execution{}, err
	}

	return model.Execution{
		Language: exec.Language,
		Code:     exec.Code,
		Tests:    tests,

		StdOut:   &execResult.StdOut,
		StdErr:   &execResult.StdErr,
		Duration: &execResult.Duration,
	}, nil
}
