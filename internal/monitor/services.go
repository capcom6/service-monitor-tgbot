package monitor

import (
	"context"
	"time"

	"github.com/capcom6/tgbot-service-monitor/internal/config"
	"github.com/capcom6/tgbot-service-monitor/internal/monitor/probes"
)

type MonitorService struct {
	Service config.Service

	p Probeer
}

func NewMonitorService(service config.Service) *MonitorService {
	svc := MonitorService{
		Service: service,
	}

	if !service.HTTPGet.IsEmpty() {
		svc.p = probes.NewHttpGet(probes.HttpGetConfig{
			TcpSocketConfig: probes.TcpSocketConfig{
				Host: service.HTTPGet.Host,
				Port: service.HTTPGet.Port,
			},
			Scheme:      service.HTTPGet.Scheme,
			Path:        service.HTTPGet.Path,
			HTTPHeaders: service.HTTPGet.HTTPHeaders.ToMap(),
		})
	} else if !service.TCPSocket.IsEmpty() {
		svc.p = probes.NewTcpSocket(probes.TcpSocketConfig{
			Host: service.TCPSocket.Host,
			Port: service.TCPSocket.Port,
		})
	}

	return &svc
}

func (s *MonitorService) Monitor(ctx context.Context, ch ProbesChannel) (err error) {
	if s.Service, err = s.Service.ApplyDefaultsAndValidate(); err != nil {
		return
	}

	go func() {
		// log.Println("Init of", s.Service.Name)

		if s.Service.InitialDelaySeconds > 0 {
			time.Sleep(time.Duration(s.Service.InitialDelaySeconds) * time.Second)
		}

		// log.Println("Start of", s.Service.Name)
		ticker := time.NewTicker(time.Duration(s.Service.PeriodSeconds) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// log.Println("Tick for", s.Service.Name)

				c, cancel := context.WithTimeout(ctx, time.Duration(s.Service.TimeoutSeconds)*time.Second)
				err = s.p.Probe(c)
				if err != nil {
					errorLog.Println("Error of", s.Service.Name, ":", err.Error())
				}
				cancel()

				select {
				case ch <- ServiceProbe{
					Id:    s.Service.Id,
					Name:  s.Service.Name,
					Error: err,
				}:
				case <-ctx.Done():
					// log.Println("Stop of", s.Service.Name, "on write")
					return
				}

			case <-ctx.Done():
				// log.Println("Stop of", s.Service.Name, "on read")
				return
			}
		}
	}()

	return
}
