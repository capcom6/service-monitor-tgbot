package storage

import "fmt"

type Storage interface {
	Load() ([]MonitoredService, error)
	Close() error
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
		if validateErr := services[i].Validate(); validateErr != nil {
			return nil, fmt.Errorf("failed to validate service %s: %w", services[i].Name, validateErr)
		}
	}

	return services, nil
}

func (s *Service) Close() error {
	if err := s.storage.Close(); err != nil {
		return fmt.Errorf("failed to close storage: %w", err)
	}

	return nil
}
