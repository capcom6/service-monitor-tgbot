package storage

type Storage interface {
	Load() ([]Service, error)
}

type StorageService struct {
	storage Storage
}

func NewStorageService(storage Storage) *StorageService {
	return &StorageService{
		storage: storage,
	}
}

func (s *StorageService) Load() ([]Service, error) {
	services, err := s.storage.Load()
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
