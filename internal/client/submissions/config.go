package submissionsclient

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"ADDR" required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("CLIENT_SUBMISSIONS", &cfg); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}
