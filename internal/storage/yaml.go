package storage

import (
	"os"

	"gopkg.in/yaml.v3"
)

type yamlRoot struct {
	Services []Service `yaml:"services"`
}

type yamlStorage struct {
	Path string
}

func (s *yamlStorage) Load() ([]Service, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}

	root := yamlRoot{}
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	return root.Services, nil
}
