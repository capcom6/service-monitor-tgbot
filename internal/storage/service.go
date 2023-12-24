package storage

import "context"

type Storage interface {
	Select(context.Context) ([]Service, error)
}

type StorageService struct {
	storage Storage
}

func NewStorageService(storage Storage) *StorageService {
	return &StorageService{
		storage: storage,
	}
}

func (s *StorageService) Select(ctx context.Context) ([]Service, error) {
	services, err := s.storage.Select(ctx)
	if err != nil {
		return nil, err
	}

	for i := range services {
		if err := services[i].Validate(); err != nil {
			return nil, err
		}
	}

	return services, nil
}
