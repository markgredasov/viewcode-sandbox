package sandbox_config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type LanguageConfig struct {
	Image string   `yaml:"image"`
	Cmd   []string `yaml:"cmd"`
}

type LanguagesConfig struct {
	Languages map[string]LanguageConfig `yaml:"languages"`
}

func newLanguagesConfig(filepath string) (LanguagesConfig, error) {
	data, err := os.ReadFile(filepath) // nolint:gosec // todo
	if err != nil {
		return LanguagesConfig{}, err
	}

	var cfg LanguagesConfig
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return LanguagesConfig{}, err
	}

	return cfg, nil
}
