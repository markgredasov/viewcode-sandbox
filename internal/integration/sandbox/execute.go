package sandbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
)

func (s *sandbox) executeTests(ctx context.Context, request executionRequest) (*executionResult, error) {
	result := &executionResult{}

	workDir, err := os.MkdirTemp("", "sandbox-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	if err = s.saveCode("main.go", workDir, request.Code); err != nil {
		return nil, err
	}

	if request.Tests != nil {
		if err = s.saveCode("main_test.go", workDir, *request.Tests); err != nil {
			return nil, err
		}
	} else {
		return nil, errs.ErrNoTestsProvided
	}

	if err = s.createGoMod(workDir); err != nil {
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

func (s *sandbox) saveCode(filename, workDir string, code string) error {
	filepath := filepath.Join(workDir, filename)

	return os.WriteFile(filepath, []byte(code), 0644)
}

func (s *sandbox) createGoMod(workDir string) error {
	moduleName := "sandbox"

	goModContent := fmt.Sprintf(`module %s

	go 1.26
	`, moduleName)

	return s.saveCode("go.mod", workDir, goModContent)
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
	image, cmd := "", []string{}

	switch language {
	case "go":
		image = "golang:latest"
		cmd = []string{"go", "test", "-v"}
	case "python3":
		image = "python:alpine"
		cmd = []string{"python", "-m", "unittest"}
	default:
		return image, cmd, fmt.Errorf("invalid language: %w", errs.ErrInvalidArgument)
	}

	return image, cmd, nil
}

func memoryToBytes(memoryInMB int64) int64 {
	var factor int64 = 1024 * 1024
	return memoryInMB * factor
}
