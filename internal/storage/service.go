package storage

import "fmt"

type Storage interface {
	Load() ([]MonitoredService, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Load() ([]MonitoredService, error) {
	services, err := s.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	for i := range services {
		if err := services[i].Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate service %s: %w", services[i].Name, err)
		}
	}

	return services, nil
}
