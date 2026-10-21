package sandbox

import (
	"github.com/moby/moby/client"

	sandboxconfig "github.com/markgredasov/viewcode-sandbox/internal/client/sandbox/config"
)

type sandbox struct {
	config       sandboxconfig.Config
	dockerClient *client.Client
}

func newSandbox(config sandboxconfig.Config) (*sandbox, error) {
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
