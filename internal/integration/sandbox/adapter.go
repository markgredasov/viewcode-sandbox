package sandbox

import (
	"context"
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

var languageModelToExecutionLanguage = map[model.Language]string{
	model.LanguageGo:         "go",
	model.LanguagePython:     "python3",
	model.LanguageJava:       "java",
	model.LanguageCpp:        "cpp",
	model.LanguageCSharp:     "csharp",
	model.LanguageJavaScript: "js",
	model.LanguageTypeScript: "ts",
	model.LanguageKotlin:     "ko",
	model.LanguageSwift:      "sw",
	model.LanguageRust:       "rust",
}

var executionLanguageToLanguageModel = map[string]model.Language{
	"go":      model.LanguageGo,
	"python3": model.LanguagePython,
	"java":    model.LanguageJava,
	"cpp":     model.LanguageCpp,
	"csharp":  model.LanguageCSharp,
	"js":      model.LanguageJavaScript,
	"ts":      model.LanguageTypeScript,
	"ko":      model.LanguageKotlin,
	"sw":      model.LanguageSwift,
	"rust":    model.LanguageRust,
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

func (a *Adapter) RunTests(
	ctx context.Context,
	exec model.Execution,
	tests model.Verification,
) (model.Execution, error) {
	language, ok := languageModelToExecutionLanguage[exec.Language]
	if !ok {
		return model.Execution{}, fmt.Errorf("invalid language: %w", errs.ErrInvalidArgument)
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

	executionResultModel := model.NewExecutionResult(execResult.StdOut, execResult.StdErr, execResult.Duration)

	return model.Execution{
		Language: exec.Language,
		Code:     exec.Code,
		Tests:    tests,

		ExecutionResult: executionResultModel,
	}, nil
}
