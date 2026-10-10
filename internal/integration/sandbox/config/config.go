package sandbox_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DockerHost           string        `envconfig:"HOST" required:"true"`
	DockerNetwork        string        `envconfig:"NETWORK" required:"true"`
	ContainerMemory      int64         `envconfig:"MEMORY_IN_MB" required:"true"`
	ContainerExecTimeout time.Duration `envconfig:"EXEC_TIMEOUT" required:"true"`
	ContainerStopTimeout time.Duration `envconfig:"STOP_TIMEOUT" required:"true"`

	ConfigLanguagesFilepath string `envconfig:"CONFIG_LANGUAGES_FILEPATH"`
	ConfigLanguages         LanguagesConfig
}

func newConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("SANDBOX_DOCKER", &config); err != nil {
		return Config{}, err
	}

	languagesConfig, err := newLanguagesConfig(config.ConfigLanguagesFilepath)
	if err != nil {
		return Config{}, err
	}
	config.ConfigLanguages = languagesConfig

	return config, nil
}

func NewConfigMust() Config {
	config, err := newConfig()
	if err != nil {
		err = fmt.Errorf("new sandbox config: %w", err)
		panic(err)
	}

	return config
}
