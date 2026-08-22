package sandbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
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

func (s *sandbox) saveCode(filename, workDir string, code string) error {
	filepath := filepath.Join(workDir, filename)

	return os.WriteFile(filepath, []byte(code), 0644)
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
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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
		//nolint:errcheck // not needed
		s.dockerClient.ContainerStop(context.Background(), containerID, client.ContainerStopOptions{})
		return 0, fmt.Errorf("execution timeout after %v", timeout)
	}
}

func (s *sandbox) getLogs(ctx context.Context, containerID string) (string, string, error) {
	stdout, err := s.dockerClient.ContainerLogs(ctx, containerID, client.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     false,
	})
	if err != nil {
		return "", "", err
	}
	defer stdout.Close()

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
	defer stderr.Close()

	stderrBytes, err := io.ReadAll(stderr)
	if err != nil {
		return "", "", err
	}

	return string(stdoutBytes), string(stderrBytes), nil
}

func (s *sandbox) cleanupContainer(ctx context.Context, containerID string) {
	//nolint:errcheck // not needed
	s.dockerClient.ContainerStop(ctx, containerID, client.ContainerStopOptions{})
	//nolint:errcheck // not needed
	s.dockerClient.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{
		Force: true,
	})
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
