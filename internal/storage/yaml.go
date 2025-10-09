package storage

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type yamlRoot struct {
	Services []MonitoredService `yaml:"services"`
}

type yamlStorage struct {
	Path string
}

func (s *yamlStorage) Load() ([]MonitoredService, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	root := new(yamlRoot)
	if err := yaml.Unmarshal(data, root); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return root.Services, nil
}
