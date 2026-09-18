package submissionsingegration

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"ADDR"`
}

func newConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("CLIENT_SUBMISSIONS", &cfg); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := newConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}
