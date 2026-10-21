package sandbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"go.uber.org/zap"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
)

func (s *sandbox) executeTests(ctx context.Context, request executionRequest) (*executionResult, error) {
	language, ok := executionLanguageToLanguageModel[request.Language]
	if !ok {
		return nil, fmt.Errorf("invalid language: %s: %w", request.Language, errs.ErrInvalidArgument)
	}

	var result *executionResult
	switch language {
	case model.LanguageGo:
		res, err := s.executeTestsLanguageGo(ctx, request)
		if err != nil {
			return result, err
		}
		result = res
	case model.LanguagePython:
		res, err := s.executeTestsLanguagePython(ctx, request)
		if err != nil {
			return result, err
		}
		result = res
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
		return nil, fmt.Errorf("language not provided yet: %w", errs.ErrNotImplemented)
	}

	return result, nil
}

func (s *sandbox) saveCode(filename, workDir, code string) error {
	filepath := filepath.Join(workDir, filename)

	return os.WriteFile(filepath, []byte(code), 0o644) //nolint:gosec // todo
}

func (s *sandbox) createContainer(ctx context.Context, workDir string, request executionRequest) (string, error) {
	memoryInBytes := memoryToBytes(s.config.ContainerMemory)
	resources := container.Resources{
		Memory:     memoryInBytes,
		MemorySwap: memoryInBytes,
	}

	image, cmd, err := s.getImageAndCmd(request.Language)
	if err != nil {
		return "", err
	}

	resp, err := s.dockerClient.ContainerCreate(
		ctx,
		client.ContainerCreateOptions{
			Config: &container.Config{
				Image:      image,
				Cmd:        cmd,
				WorkingDir: "/app",
			},
			HostConfig: &container.HostConfig{
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: workDir,
						Target: "/app",
					},
				},
				Resources:   resources,
				AutoRemove:  false,
				NetworkMode: container.NetworkMode("none"),
			},
			Name: fmt.Sprintf("sandbox-%d", time.Now().UnixNano()),
		},
	)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (s *sandbox) startContainer(ctx context.Context, containerID string) error {
	_, err := s.dockerClient.ContainerStart(ctx, containerID, client.ContainerStartOptions{})
	return err
}

func (s *sandbox) waitContainer(ctx context.Context, containerID string, timeout time.Duration) (int64, error) {
	log := logger.FromContext(ctx)

	waitCtx, cancelWaitCtx := context.WithTimeout(ctx, timeout)
	defer cancelWaitCtx()

	statusCh := s.dockerClient.ContainerWait(
		waitCtx,
		containerID,
		client.ContainerWaitOptions{
			Condition: container.WaitConditionNotRunning,
		},
	)

	select {
	case status := <-statusCh.Result:
		if status.Error != nil {
			return 0, fmt.Errorf("container error: %s", status.Error.Message)
		}
		return status.StatusCode, nil
	case <-waitCtx.Done():
		stopCtx, cancelStopCtx := context.WithTimeout(context.WithoutCancel(ctx), s.config.ContainerStopTimeout)
		defer cancelStopCtx()

		_, err := s.dockerClient.ContainerStop(stopCtx, containerID, client.ContainerStopOptions{})
		if err != nil {
			log.Error("container stop", zap.Error(err))
		}
		return 0, errs.TimeoutExceeded(fmt.Sprintf("execution timeout after %v", timeout))
	}
}

func (s *sandbox) getLogs(ctx context.Context, containerID string) (string, string, error) {
	log := logger.FromContext(ctx)
	stdout, err := s.dockerClient.ContainerLogs(ctx, containerID, client.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     false,
	})
	if err != nil {
		return "", "", err
	}

	defer func() {
		err = stdout.Close()
		if err != nil {
			log.Error("close stdout", zap.Error(err))
		}
	}()

	stdoutBytes, err := io.ReadAll(stdout)
	if err != nil {
		return "", "", err
	}

	stderr, err := s.dockerClient.ContainerLogs(ctx, containerID, client.ContainerLogsOptions{
		ShowStderr: true,
		Follow:     false,
	})
	if err != nil {
		return "", "", err
	}
	defer func() {
		err = stderr.Close()
		if err != nil {
			log.Error("close stderr", zap.Error(err))
		}
	}()

	stderrBytes, err := io.ReadAll(stderr)
	if err != nil {
		return "", "", err
	}

	return string(stdoutBytes), string(stderrBytes), nil
}

func (s *sandbox) cleanupContainer(ctx context.Context, containerID string) {
	log := logger.FromContext(ctx)
	_, err := s.dockerClient.ContainerStop(ctx, containerID, client.ContainerStopOptions{})
	if err != nil {
		log.Error("container stop", zap.Error(err))
	}

	_, err = s.dockerClient.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Error("container remove", zap.Error(err))
	}
}

func (s *sandbox) getImageAndCmd(language string) (string, []string, error) {
	languageModel, ok := executionLanguageToLanguageModel[language]
	if !ok {
		return "", []string{}, fmt.Errorf("invalid language: %s: %w", language, errs.ErrInvalidArgument)
	}

	languageConfig, ok := s.config.ConfigLanguages.Languages[languageModel.String()]
	if !ok {
		return "", []string{}, fmt.Errorf("invalid language: %s: %w", language, errs.ErrInvalidArgument)
	}

	return languageConfig.Image, languageConfig.Cmd, nil
}

func memoryToBytes(memoryInMB int64) int64 {
	var factor int64 = 1024 * 1024
	return memoryInMB * factor
}
