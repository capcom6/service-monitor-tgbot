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
		mu:     sync.RWMutex{},
	}, nil
}

// initializeProbes creates monitoring tasks for all services and returns any initialization error.
func (m *Service) initializeProbes(ctx context.Context) error {
	for i, cfg := range m.services {
		mon, err := newTask(newTaskConfig(cfg))
		if err != nil {
			return NewProbeInitializationError(cfg.ID, cfg.Name, err)
		}

		monCh, err := mon.Monitor(ctx)
		if err != nil {
			return NewProbeInitializationError(cfg.ID, cfg.Name, err)
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
		go func(i int, ch ProbesChannel) {
			defer wg.Done()
			m.startProbe(ctx, i, ch, updCh)
		}(i, ch)
	}

	wg.Wait()
	m.logger.Info("Monitor service stopped")
}

func (m *Service) startProbe(ctx context.Context, id int, ch ProbesChannel, updCh UpdatesChannel) {
	m.logger.Info("Starting probe", zap.String("id", m.services[id].ID))

	for {
		select {
		case probe, ok := <-ch:
			if !ok {
				m.logger.Warn("Probe channel closed", zap.String("id", m.services[id].ID))
				return
			}
			if update := m.updateState(id, probe); update != nil {
				select {
				case updCh <- *update:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			m.logger.Info("Stopping probe", zap.String("id", m.services[id].ID))
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
		current.Probes += delta
	} else {
		current.Probes = delta
	}

	cooldown := time.Duration(service.AlertCooldownSeconds) * time.Second
	inCooldown := cooldown > 0 && time.Since(current.LastAlertedAt) < cooldown

	var upd *ServiceStatus
	if !current.Online && current.Probes >= int(service.SuccessThreshold) {
		upd = m.handleTransitionToOnline(&current, service, inCooldown)
	} else if current.Online && current.Probes <= -int(service.FailureThreshold) {
		upd = m.handleTransitionToOffline(&current, service, probe, inCooldown)
	}

	current.Timestamp = time.Now()

	m.mu.Lock()
	m.states[id] = current
	m.mu.Unlock()

	return upd
}

func (m *Service) handleTransitionToOnline(
	current *state,
	service storage.MonitoredService,
	inCooldown bool,
) *ServiceStatus {
	if inCooldown {
		current.Probes = int(service.SuccessThreshold) - 1
		return nil
	}

	current.Online = true
	current.Error = nil
	current.ChangedAt = time.Now()
	current.LastAlertedAt = time.Now()

	return &ServiceStatus{
		ID:        service.ID,
		Name:      service.Name,
		State:     ServiceStateOnline,
		Error:     nil,
		ChangedAt: current.ChangedAt,
	}
}

func (m *Service) handleTransitionToOffline(
	current *state,
	service storage.MonitoredService,
	probe error,
	inCooldown bool,
) *ServiceStatus {
	if inCooldown {
		current.Probes = -(int(service.FailureThreshold) - 1)
		return nil
	}

	current.Online = false
	current.Error = probe
	current.ChangedAt = time.Now()
	current.LastAlertedAt = time.Now()

	return &ServiceStatus{
		ID:        service.ID,
		Name:      service.Name,
		State:     ServiceStateOffline,
		Error:     probe,
		ChangedAt: current.ChangedAt,
	}
}

// GetCurrentStatuses returns the current status of all services.
func (m *Service) GetCurrentStatuses() []ServiceStatus {
	statuses := make([]ServiceStatus, len(m.services))

	m.mu.RLock()
	defer m.mu.RUnlock()

	for i, service := range m.services {
		state := m.states[i]
		statuses[i] = ServiceStatus{
			ID:        service.ID,
			Name:      service.Name,
			State:     state.State(),
			Error:     state.Error,
			ChangedAt: state.ChangedAt,
		}
	}

	return statuses
}
