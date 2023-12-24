package monitor

import (
	"context"
	"fmt"
	"sync"

	"github.com/capcom6/service-monitor-tgbot/internal/monitor/probes"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
	"github.com/capcom6/service-monitor-tgbot/pkg/collections"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"monitor",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("monitor")
	}),
	fx.Provide(NewMonitorModule),
)

type MonitorModuleParams struct {
	fx.In

	Storage storage.Storage
	Logger  *zap.Logger
}

type MonitorModule struct {
	Storage  storage.Storage
	Services []storage.Service
	Logger   *zap.Logger

	probes []ProbesChannel
	states []state
}

func NewMonitorModule(params MonitorModuleParams) (*MonitorModule, error) {
	services, err := params.Storage.Select(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	return &MonitorModule{
		Storage:  params.Storage,
		Services: services,
		Logger:   params.Logger,
	}, nil
}

func (m *MonitorModule) Monitor(ctx context.Context) (UpdatesChannel, error) {
	m.probes = make([]ProbesChannel, len(m.Services))
	m.states = make([]state, len(m.Services))

	for i, cfg := range m.Services {
		mon := NewMonitorService(ServiceMonitorConfig{
			HttpGet: probes.HttpGetConfig{
				TcpSocketConfig: probes.TcpSocketConfig{
					Host: cfg.HTTPGet.TCPSocket.Host,
					Port: cfg.HTTPGet.TCPSocket.Port,
				},
				Scheme:      cfg.HTTPGet.Scheme,
				Path:        cfg.HTTPGet.Path,
				HTTPHeaders: collections.GroupBy(cfg.HTTPGet.HTTPHeaders, func(h storage.HTTPHeader) (string, string) { return h.Name, h.Value }),
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
				m.Logger.Info("Starting probe", zap.Int("id", i))
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
						m.Logger.Info("Stopping probe", zap.Int("id", i))
						return
					}
				}
			}(i, ch)
		}

		wg.Wait()
		<-ctx.Done()
		m.Logger.Info("Monitor service stopped")
	}()

	return updCh, nil
}

func (m *MonitorModule) updateState(id int, probe error) *ServiceStatus {
	service := m.Services[id]
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
		upd = &ServiceStatus{
			Id:    service.Id,
			Name:  service.Name,
			State: ServiceOnline,
			Error: nil,
		}
	} else if current.Online && current.Probes == -int(service.FailureThreshold) {
		current.Online = false
		upd = &ServiceStatus{
			Id:    service.Id,
			Name:  service.Name,
			State: ServiceOffline,
			Error: probe,
		}
	}

	m.states[id] = current

	return upd
}
