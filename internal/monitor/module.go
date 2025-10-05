package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/capcom6/go-infra-fx/fxutil"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor/probes"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"monitor",
		fxutil.WithNamedLogger("monitor"),
		fx.Provide(NewMonitorModule),
	)
}

type MonitorModuleParams struct {
	fx.In

	Storage storage.Storage
	Logger  *zap.Logger
}

type MonitorModule struct {
	storage  storage.Storage
	services []storage.Service
	logger   *zap.Logger

	probes []ProbesChannel
	states []state
	mu     sync.RWMutex
}

func NewMonitorModule(params MonitorModuleParams) (*MonitorModule, error) {
	services, err := params.Storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	return &MonitorModule{
		storage:  params.Storage,
		services: services,
		logger:   params.Logger,
	}, nil
}

func (m *MonitorModule) Monitor(ctx context.Context) (UpdatesChannel, error) {
	m.probes = make([]ProbesChannel, len(m.services))
	m.states = make([]state, len(m.services))

	for i, cfg := range m.services {
		mon := NewMonitorService(ServiceMonitorConfig{
			HttpGet: probes.HttpGetConfig{
				TcpSocketConfig: probes.TcpSocketConfig{
					Host: cfg.HTTPGet.Host,
					Port: cfg.HTTPGet.Port,
				},
				Scheme:      cfg.HTTPGet.Scheme,
				Path:        cfg.HTTPGet.Path,
				HTTPHeaders: lo.GroupByMap(cfg.HTTPGet.HTTPHeaders, func(h storage.HTTPHeader) (string, string) { return h.Name, h.Value }),
			},
			TcpSocket: probes.TcpSocketConfig{
				Host: cfg.TCPSocket.Host,
				Port: cfg.TCPSocket.Port,
			},
			InitialDelaySeconds: uint16(cfg.InitialDelaySeconds),
			PeriodSeconds:       cfg.PeriodSeconds,
			TimeoutSeconds:      cfg.TimeoutSeconds,
		})
		monCh, err := mon.Monitor(ctx)
		if err != nil {
			return nil, err
		}
		m.probes[i] = monCh
	}

	updCh := make(UpdatesChannel)
	go func() {
		defer close(updCh)

		wg := sync.WaitGroup{}
		wg.Add(len(m.probes))

		for i, ch := range m.probes {
			go func(i int, ch ProbesChannel) {
				m.logger.Info("Starting probe", zap.Int("id", i))
				defer wg.Done()
				for {
					select {
					case probe := <-ch:
						if update := m.updateState(i, probe); update != nil {
							select {
							case updCh <- *update:
							case <-ctx.Done():
								return
							}
						}
					case <-ctx.Done():
						m.logger.Info("Stopping probe", zap.Int("id", i))
						return
					}
				}
			}(i, ch)
		}

		wg.Wait()
		<-ctx.Done()
		m.logger.Info("Monitor service stopped")
	}()

	return updCh, nil
}

func (m *MonitorModule) updateState(id int, probe error) *ServiceStatus {
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
			Id:    service.Id,
			Name:  service.Name,
			State: ServiceStateOnline,
			Error: nil,
		}
	} else if current.Online && current.Probes == -int(service.FailureThreshold) {
		current.Online = false
		current.Error = probe

		upd = &ServiceStatus{
			Id:    service.Id,
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
func (m *MonitorModule) GetCurrentStatuses() []ServiceStatus {
	statuses := make([]ServiceStatus, len(m.services))

	m.mu.RLock()
	defer m.mu.RUnlock()

	for i, service := range m.services {
		state := m.states[i]
		statuses[i] = ServiceStatus{
			Id:    service.Id,
			Name:  service.Name,
			State: state.State(),
			Error: state.Error,
		}
	}

	return statuses
}
