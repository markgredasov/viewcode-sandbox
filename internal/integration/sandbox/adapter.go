package sandbox

import (
	"context"
	"fmt"
	"time"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

var languageToExecutionLanguage = map[model.Language]string{
	model.LanguageGo:     "go",
	model.LanguagePython: "python3",
	model.LanguageRust:   "rust",
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

func (a *Adapter) RunCode(ctx context.Context, exec model.Execution) (model.Execution, error) {
	language, ok := languageToExecutionLanguage[exec.Language]
	if !ok {
		return model.Execution{}, errs.ErrInvalidLanguage
	}

	req := executionRequest{
		Language: language,
		Code:     exec.Code.String(),
		Timeout:  10 * time.Second,
	}

	execResult, err := a.sandbox.execute(ctx, req)
	if err != nil {
		return model.Execution{}, err
	}

	return model.Execution{
		Language:         exec.Language,
		Code:             exec.Code,
		InlineExpression: exec.InlineExpression,
		Arguments:        exec.Arguments,
		Environment:      exec.Environment,

		StdOut:   &execResult.StdOut,
		StdErr:   &execResult.StdErr,
		Duration: &execResult.Duration,
	}, nil
}
