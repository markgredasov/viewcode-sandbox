package sandbox

import (
	"context"
	"fmt"
	"os"
	"time"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"go.uber.org/zap"
)

func (s *sandbox) executeTestsLanguagePython(ctx context.Context, request executionRequest) (*executionResult, error) {
	log := logger.FromContext(ctx)

	result := &executionResult{}

	workDir, err := os.MkdirTemp("", "sandbox-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() {
		err = os.RemoveAll(workDir)
		if err != nil {
			log.Error("remove all", zap.Error(err))
		}
	}()

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
