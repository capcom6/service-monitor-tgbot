package storage

import (
	"context"
	"net/url"
	"sync"

	"github.com/capcom6/service-monitor-tgbot/pkg/collections"
)

type memoryStorage struct {
	services map[string]Service
	mux      sync.RWMutex
}

func newMemoryStorage(u *url.URL) (*memoryStorage, error) {
	return &memoryStorage{
		services: make(map[string]Service),
	}, nil
}

func (s *memoryStorage) Select(ctx context.Context) ([]Service, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return collections.Values(s.services), nil
}
