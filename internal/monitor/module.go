package monitor

import (
	"context"
	"fmt"
	"sync"

	"github.com/capcom6/tgbot-service-monitor/internal/config"
	"github.com/capcom6/tgbot-service-monitor/internal/monitor/probes"
)

type MonitorModule struct {
	Services []config.Service

	probes []ProbesChannel
	states []int
}

func NewMonitorModule(services []config.Service) *MonitorModule {
	return &MonitorModule{
		Services: services,
	}
}

func (m *MonitorModule) Monitor(ctx context.Context) (UpdatesChannel, error) {
	m.probes = make([]ProbesChannel, len(m.Services))
	m.states = make([]int, len(m.Services))

	for i, s := range m.Services {
		cfg, err := s.ApplyDefaultsAndValidate()
		if err != nil {
			return nil, fmt.Errorf("invalid config for %s: %w", s.Name, err)
		}
		mon := NewMonitorService(ServiceMonitorConfig{
			HttpGet: probes.HttpGetConfig{
				TcpSocketConfig: probes.TcpSocketConfig{
					Host: cfg.HTTPGet.TCPSocket.Host,
					Port: cfg.HTTPGet.TCPSocket.Port,
				},
				Scheme:      cfg.HTTPGet.Scheme,
				Path:        cfg.HTTPGet.Path,
				HTTPHeaders: cfg.HTTPGet.HTTPHeaders.ToMap(),
			},
			TcpSocket: probes.TcpSocketConfig{
				Host: cfg.TCPSocket.Host,
				Port: cfg.TCPSocket.Port,
			},
			InitialDelaySeconds: cfg.InitialDelaySeconds,
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
						log.Println("Probe", i, "stopped")
						return
					}
				}
			}(i, ch)
		}

		wg.Wait()
		<-ctx.Done()
		log.Println("Monitor service stopped")
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

	if current > 0 && delta > 0 {
		// защита от переполнения
		if current <= service.SuccessThreshold {
			current += delta
		}
	} else if current < 0 && delta < 0 {
		if current >= -service.FailureThreshold {
			current += delta
		}
	} else {
		current = delta
	}

	var upd *ServiceStatus
	if current == service.SuccessThreshold {
		upd = &ServiceStatus{
			Id:    service.Id,
			Name:  service.Name,
			State: ServiceOnline,
			Error: nil,
		}
	} else if current == -service.FailureThreshold {
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
