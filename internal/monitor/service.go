package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/storage"
	"go.uber.org/zap"
)

type Service struct {
	storage  storage.Storage
	services []storage.MonitoredService
	logger   *zap.Logger

	probes []ProbesChannel
	states []state
	mu     sync.RWMutex
}

func NewService(storage storage.Storage, logger *zap.Logger) (*Service, error) {
	services, err := storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	return &Service{
		storage:  storage,
		services: services,
		logger:   logger,

		probes: make([]ProbesChannel, len(services)),
		states: make([]state, len(services)),
	}, nil
}

// initializeProbes creates monitoring tasks for all services and returns any initialization error
func (m *Service) initializeProbes(ctx context.Context) error {
	for i, cfg := range m.services {
		mon, err := newTask(newTaskConfig(cfg))
		if err != nil {
			m.logger.Error("failed to initialize probe for service", zap.Int("id", i), zap.Error(err))
			continue
		}

		monCh, err := mon.Monitor(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize probe for service %d: %w", i, err)
		}
		m.probes[i] = monCh
	}

	return nil
}

func (m *Service) Monitor(ctx context.Context) (UpdatesChannel, error) {
	if err := m.initializeProbes(ctx); err != nil {
		return nil, err
	}

	updCh := make(UpdatesChannel)
	go m.startMonitoring(ctx, updCh)

	return updCh, nil
}

func (m *Service) startMonitoring(ctx context.Context, updCh UpdatesChannel) {
	defer close(updCh)

	wg := sync.WaitGroup{}
	wg.Add(len(m.probes))

	for i, ch := range m.probes {
		go m.startProbe(ctx, i, ch, &wg, updCh)
	}

	wg.Wait()
	<-ctx.Done()
	m.logger.Info("Monitor service stopped")
}

func (m *Service) startProbe(ctx context.Context, id int, ch ProbesChannel, wg *sync.WaitGroup, updCh UpdatesChannel) {
	defer wg.Done()
	m.logger.Info("Starting probe", zap.Int("id", id))

	for {
		select {
		case probe := <-ch:
			if update := m.updateState(id, probe); update != nil {
				select {
				case updCh <- *update:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			m.logger.Info("Stopping probe", zap.Int("id", id))
			return
		}
	}
}

func (m *Service) updateState(id int, probe error) *ServiceStatus {
	service := m.services[id]
	current := m.states[id]

	delta := 1
	if probe != nil {
		delta = -1
	}

	if (current.Probes > 0 && delta > 0) ||
		(current.Probes < 0 && delta < 0) {
		// если знак совпадает, то продолжаем
		current.Probes += delta
	} else {
		current.Probes = delta
	}

	var upd *ServiceStatus
	if !current.Online && current.Probes == int(service.SuccessThreshold) {
		current.Online = true
		current.Error = nil

		upd = &ServiceStatus{
			ID:    service.ID,
			Name:  service.Name,
			State: ServiceStateOnline,
			Error: nil,
		}
	} else if current.Online && current.Probes == -int(service.FailureThreshold) {
		current.Online = false
		current.Error = probe

		upd = &ServiceStatus{
			ID:    service.ID,
			Name:  service.Name,
			State: ServiceStateOffline,
			Error: probe,
		}
	}

	current.Timestamp = time.Now()

	m.mu.Lock()
	m.states[id] = current
	m.mu.Unlock()

	return upd
}

// GetCurrentStatuses returns the current status of all services
func (m *Service) GetCurrentStatuses() []ServiceStatus {
	statuses := make([]ServiceStatus, len(m.services))

	m.mu.RLock()
	defer m.mu.RUnlock()

	for i, service := range m.services {
		state := m.states[i]
		statuses[i] = ServiceStatus{
			ID:    service.ID,
			Name:  service.Name,
			State: state.State(),
			Error: state.Error,
		}
	}

	return statuses
}
