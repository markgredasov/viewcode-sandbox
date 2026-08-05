package sandbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
)

type sandbox struct {
	config       sandbox_config.Config
	dockerClient *client.Client
}

func newSandbox(config sandbox_config.Config) (*sandbox, error) {
	client, err := client.New(
		client.FromEnv,
		client.WithHost(config.DockerHost),
	)
	if err != nil {
		return &sandbox{}, err
	}

	return &sandbox{
		dockerClient: client,
		config:       config,
	}, nil
}

func (s *sandbox) execute(ctx context.Context, request executionRequest) (*executionResult, error) {
	startTime := time.Now()
	result := &executionResult{}

	workDir, err := os.MkdirTemp("", "sandbox-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	if err := s.saveCode(workDir, request); err != nil {
		return nil, err
	}

	containerID, err := s.createContainer(ctx, workDir, request)
	if err != nil {
		return nil, err
	}
	defer s.cleanupContainer(ctx, containerID)

	if err := s.startContainer(ctx, containerID); err != nil {
		return nil, err
	}

	exitCode, err := s.waitContainer(ctx, containerID, request.Timeout)
	if err != nil {
		return nil, err
	}
	result.ExitCode = exitCode

	stdout, stderr, err := s.getLogs(ctx, containerID)
	if err != nil {
		return nil, err
	}
	result.StdOut = stdout
	result.StdErr = stderr
	result.Duration = time.Since(startTime)

	return result, nil
}

func (s *sandbox) saveCode(workDir string, request executionRequest) error {
	filename := "main.go"
	filepath := filepath.Join(workDir, filename)

	return os.WriteFile(filepath, []byte(request.Code), 0644)
}

func (s *sandbox) createContainer(ctx context.Context, workDir string, request executionRequest) (string, error) {
	memoryInBytes := memoryToBytes(s.config.DockerContainerMemory)
	resources := container.Resources{
		Memory:     memoryInBytes,
		MemorySwap: memoryInBytes,
	}

	image, cmd, err := s.getImageAndCmd(ctx, request.Language)
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

	statusCh := s.dockerClient.ContainerWait(waitCtx, containerID, client.ContainerWaitOptions{Condition: container.WaitConditionNotRunning})

	select {
	case status := <-statusCh.Result:
		if status.Error != nil {
			return 0, fmt.Errorf("container error: %s", status.Error.Message)
		}
		return status.StatusCode, nil
	case <-waitCtx.Done():
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
	s.dockerClient.ContainerStop(ctx, containerID, client.ContainerStopOptions{})
	s.dockerClient.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{
		Force: true,
	})
}

func (s *sandbox) getImageAndCmd(ctx context.Context, language string) (string, []string, error) {
	image, cmd := "", []string{}

	switch language {
	case "go":
		image = "golang:latest"
		cmd = []string{"go", "run", "main.go"}
	default:
		return image, cmd, errs.ErrInvalidLanguage
	}

	return image, cmd, nil
}

func memoryToBytes(memoryInMB int64) int64 {
	return memoryInMB * 1024 * 1024
}
