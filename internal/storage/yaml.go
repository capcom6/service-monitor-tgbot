package storage

import (
	"fmt"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type yamlRoot struct {
	Services []MonitoredService `yaml:"services"`
}

type YamlStorage struct {
	Path string
}

func newYamlStorage(u *url.URL) (Storage, error) {
	var path string
	switch {
	case u.Opaque != "":
		path = u.Opaque
	case u.Host == "localhost":
		path = u.Path
	case u.Host != "":
		path = u.Host + u.Path
	default:
		path = u.Path
	}

	return &YamlStorage{
		Path: path,
	}, nil
}

func (s *YamlStorage) Load() ([]MonitoredService, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	root := new(yamlRoot)
	if yamlErr := yaml.Unmarshal(data, root); yamlErr != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", yamlErr)
	}

	return root.Services, nil
}

func (s *YamlStorage) Close() error {
	return nil
}
