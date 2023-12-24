package storage

import (
	"context"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type yamlRoot struct {
	Services []Service `yaml:"services"`
}

type yamlStorage struct {
	Path string
}

func newYamlStorage(u *url.URL) (*yamlStorage, error) {
	path := u.Path
	if u.Host == "." {
		path = "./" + path
	}

	return &yamlStorage{
		Path: path,
	}, nil
}

func (s *yamlStorage) Select(ctx context.Context) ([]Service, error) {
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
