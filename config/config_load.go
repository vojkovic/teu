package config

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*TeuConfig, error) {
	var _ io.Reader = (*os.File)(nil)

	fileContent, err := os.ReadFile(filepath.Clean(path + "/teu.yml"))
	if err != nil {
		return nil, err
	}

	var cfg TeuConfig
	if err := yaml.Unmarshal(fileContent, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(path); err != nil {
		return nil, err
	}

	return &cfg, nil
}