package monitor

import (
	"context"
	"strconv"

	"github.com/capcom6/tgbot-service-monitor/internal/config"
)

type MonitorModule struct {
	Services []config.Service

	monitors []*MonitorService
	states   []ServiceState
}

func NewMonitorModule(services []config.Service) *MonitorModule {
	return &MonitorModule{
		Services: services,
	}
}

func (m *MonitorModule) Monitor(ctx context.Context) (UpdatesChannel, error) {
	m.monitors = make([]*MonitorService, len(m.Services))
	m.states = make([]ServiceState, len(m.Services))

	probes := make(ProbesChannel)
	for i, s := range m.Services {
		// хак, чтобы обеспечить уникальность и доступ по индексу
		// по-хорошему надо использовать внутреннюю структуру, не связанную с конфигом
		s.Id = strconv.Itoa(i)

		m.monitors[i] = NewMonitorService(s)
		if err := m.monitors[i].Monitor(ctx, probes); err != nil {
			close(probes)
			return nil, err
		}
	}

	ch := make(UpdatesChannel)
	go func() {
		defer close(probes)
		defer close(ch)
		for {
			select {
			case probe := <-probes:
				if update := m.updateState(probe); update != nil {
					ch <- *update
				}
			case <-ctx.Done():
				log.Println("Monitor service stopped")
				return
			}
		}
	}()

	return ch, nil
}

func (m *MonitorModule) updateState(probe ServiceProbe) *ServiceStatus {
	id, _ := strconv.Atoi(probe.Id)
	current := m.states[id]

	if current != ServiceOnline && probe.Error == nil {
		m.states[id] = ServiceOnline
		return &ServiceStatus{
			Id:    m.Services[id].Id,
			Name:  m.Services[id].Name,
			State: ServiceOnline,
			Error: nil,
		}
	}
	if current != ServiceOffline && probe.Error != nil {
		m.states[id] = ServiceOffline
		return &ServiceStatus{
			Id:    m.Services[id].Id,
			Name:  m.Services[id].Name,
			State: ServiceOffline,
			Error: probe.Error,
		}
	}

	return nil
}
