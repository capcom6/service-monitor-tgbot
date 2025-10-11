package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/monitor/probes"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
	"github.com/samber/lo"
)

type taskConfig struct {
	HTTPGet             probes.HTTPGetConfig
	TCPSocket           probes.TCPSocketConfig
	InitialDelaySeconds uint16
	PeriodSeconds       uint16
	TimeoutSeconds      uint16
}

func newTaskConfig(service storage.MonitoredService) taskConfig {
	return taskConfig{
		HTTPGet: probes.HTTPGetConfig{
			TCPSocketConfig: probes.TCPSocketConfig{
				Host: service.HTTPGet.Host,
				Port: service.HTTPGet.Port,
			},
			Scheme: service.HTTPGet.Scheme,
			Path:   service.HTTPGet.Path,
			HTTPHeaders: lo.GroupByMap(
				service.HTTPGet.HTTPHeaders,
				func(h storage.HTTPHeader) (string, string) { return h.Name, h.Value },
			),
		},
		TCPSocket: probes.TCPSocketConfig{
			Host: service.TCPSocket.Host,
			Port: service.TCPSocket.Port,
		},
		InitialDelaySeconds: service.InitialDelaySeconds(),
		PeriodSeconds:       service.PeriodSeconds,
		TimeoutSeconds:      service.TimeoutSeconds,
	}
}

type task struct {
	config taskConfig

	p Probeer
}

func newTask(config taskConfig) (*task, error) {
	var p Probeer
	switch {
	case config.HTTPGet.Host != "":
		p = probes.NewHTTPGet(config.HTTPGet)
	case config.TCPSocket.Host != "":
		p = probes.NewTCPSocket(config.TCPSocket)
	default:
		return nil, fmt.Errorf("%w: no probe configured", ErrInvalidConfig)
	}

	svc := task{
		config: config,
		p:      p,
	}

	return &svc, nil
}

func (s *task) Monitor(ctx context.Context) (ProbesChannel, error) {
	if s.p == nil {
		return nil, fmt.Errorf("%w: no probeer configured", ErrInvalidConfig)
	}

	ch := make(ProbesChannel)
	go func() {
		defer close(ch)

		if s.config.InitialDelaySeconds > 0 {
			timer := time.NewTimer(time.Duration(s.config.InitialDelaySeconds) * time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}

		ticker := time.NewTicker(time.Duration(s.config.PeriodSeconds) * time.Second)
		defer ticker.Stop()

		for {
			c, cancel := context.WithTimeout(ctx, time.Duration(s.config.TimeoutSeconds)*time.Second)
			err := s.p.Probe(c)
			cancel()

			select {
			case ch <- err:
			case <-ctx.Done():
				return
			}

			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}
