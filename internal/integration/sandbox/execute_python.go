package sandbox

import (
	"context"
	"fmt"
	"os"
	"time"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

func (s *sandbox) executeTestsLanguagePython(ctx context.Context, request executionRequest) (*executionResult, error) {
	result := &executionResult{}

	workDir, err := os.MkdirTemp("", "sandbox-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	if err = s.saveCode("main.py", workDir, request.Code); err != nil {
		return nil, err
	}

	if request.Tests == nil {
		return nil, errs.NotFound("no tests found")
	}

	if err = s.saveCode("test_main.py", workDir, *request.Tests); err != nil {
		return nil, err
	}

	containerID, err := s.createContainer(ctx, workDir, request)
	if err != nil {
		return nil, err
	}
	defer s.cleanupContainer(ctx, containerID)

	startTime := time.Now()
	if err = s.startContainer(ctx, containerID); err != nil {
		return nil, err
	}

	exitCode, err := s.waitContainer(ctx, containerID, request.Timeout)
	if err != nil {
		return nil, err
	}
	result.ExitCode = exitCode
	result.Duration = time.Since(startTime)

	stdout, stderr, err := s.getLogs(ctx, containerID)
	if err != nil {
		return nil, err
	}
	result.StdOut = stdout
	result.StdErr = stderr

	return result, nil
}
