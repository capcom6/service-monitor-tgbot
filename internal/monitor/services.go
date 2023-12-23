package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/monitor/probes"
)

type ServiceMonitorConfig struct {
	HttpGet             probes.HttpGetConfig
	TcpSocket           probes.TcpSocketConfig
	InitialDelaySeconds uint16
	PeriodSeconds       uint16
	TimeoutSeconds      uint16
}

type MonitorService struct {
	Config ServiceMonitorConfig

	p Probeer
}

func NewMonitorService(config ServiceMonitorConfig) *MonitorService {
	svc := MonitorService{
		Config: config,
	}

	if config.HttpGet.Host != "" {
		svc.p = probes.NewHttpGet(config.HttpGet)
	} else if config.TcpSocket.Host != "" {
		svc.p = probes.NewTcpSocket(config.TcpSocket)
	}

	return &svc
}

func (s *MonitorService) Monitor(ctx context.Context) (ch ProbesChannel, err error) {
	if s.p == nil {
		return nil, fmt.Errorf("got invalid config, no probeer")
	}

	ch = make(ProbesChannel)
	go func() {
		defer close(ch)

		if s.Config.InitialDelaySeconds > 0 {
			time.Sleep(time.Duration(s.Config.InitialDelaySeconds) * time.Second)
		}

		ticker := time.NewTicker(time.Duration(s.Config.PeriodSeconds) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c, cancel := context.WithTimeout(ctx, time.Duration(s.Config.TimeoutSeconds)*time.Second)
				err = s.p.Probe(c)
				cancel()

				select {
				case ch <- err:
				case <-ctx.Done():
					return
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return
}
