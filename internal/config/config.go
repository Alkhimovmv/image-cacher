package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Load(path string) (Config, error) {
	data, error := os.ReadFile(path)
	if error != nil {
		return Config{}, error
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
