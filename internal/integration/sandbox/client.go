package sandbox

import (
	"github.com/moby/moby/client"

	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
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
