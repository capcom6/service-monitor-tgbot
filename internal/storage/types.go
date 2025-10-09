package storage

import (
	"fmt"
	"math/rand/v2"
)

// MonitoredService is a service that needs to be monitored.
type MonitoredService struct {
	// ID is the unique identifier of the service.
	ID string `yaml:"id" validate:"required"`
	// Name is the human-readable name of the service.
	Name string `yaml:"name" validate:"required"`
	// InitialDelaySecondsRaw is the number of seconds to wait before starting the monitoring.
	InitialDelaySecondsRaw int16 `yaml:"initialDelaySeconds"`
	// PeriodSeconds is the number of seconds between each monitoring attempt.
	PeriodSeconds uint16 `yaml:"periodSeconds" validate:"gt=0"`
	// TimeoutSeconds is the number of seconds to wait for a response from the service.
	TimeoutSeconds uint16 `yaml:"timeoutSeconds" validate:"gt=0"`
	// SuccessThreshold is the number of successful monitoring attempts needed to consider the service as "online".
	SuccessThreshold uint8 `yaml:"successThreshold" validate:"gt=0"`
	// FailureThreshold is the number of failed monitoring attempts needed to consider the service as "offline".
	FailureThreshold uint8 `yaml:"failureThreshold" validate:"gt=0"`
	// HTTPGet is the HTTP request to send to the service.
	HTTPGet HTTPGet `yaml:"httpGet,omitempty" validate:"required_without=TCPSocket"`
	// TCPSocket is the TCP socket to connect to the service.
	TCPSocket TCPSocket `yaml:"tcpSocket,omitempty" validate:"required_without=HTTPGet"`
}

func (s *MonitoredService) InitialDelaySeconds() uint16 {
	if s.InitialDelaySecondsRaw < 0 {
		return uint16(rand.IntN(int(s.PeriodSeconds) + 1)) //nolint:gosec //weak random is enough
	}

	return uint16(s.InitialDelaySecondsRaw)
}

func (s *MonitoredService) Validate() error {
	if s.ID == "" {
		s.ID = s.Name
	}
	if s.Name == "" {
		return fmt.Errorf("%w: name is empty", ErrValidation)
	}
	if s.PeriodSeconds == 0 {
		s.PeriodSeconds = 10
	}
	if s.TimeoutSeconds == 0 {
		s.TimeoutSeconds = 1
	}

	if s.SuccessThreshold == 0 {
		s.SuccessThreshold = 1
	}
	if s.FailureThreshold == 0 {
		s.FailureThreshold = 3
	}

	if !s.HTTPGet.IsEmpty() {
		return s.HTTPGet.Validate()
	}
	if !s.TCPSocket.IsEmpty() {
		return s.TCPSocket.Validate()
	}

	return fmt.Errorf("%w: one of httpGet or tcpSocket must be filled", ErrValidation)
}

type HTTPGet struct {
	TCPSocket   `yaml:",inline"`
	Scheme      string      `yaml:"scheme"`
	Path        string      `yaml:"path"`
	HTTPHeaders HTTPHeaders `yaml:"httpHeaders"`
}

func (s *HTTPGet) Validate() error {
	if s.Scheme == "" {
		s.Scheme = "http"
	}
	if s.Path == "" {
		s.Path = "/"
	}

	if s.Port == 0 {
		switch s.Scheme {
		case "http":
			s.Port = 80
		case "https":
			s.Port = 443
		default:
			return fmt.Errorf("%w: invalid scheme %q", ErrValidation, s.Scheme)
		}
	}

	return s.TCPSocket.Validate()
}

type HTTPHeaders []HTTPHeader

type HTTPHeader struct {
	Name  string `yaml:"name" validate:"required"`
	Value string `yaml:"value" validate:"required"`
}

type TCPSocket struct {
	Host string `yaml:"host" validate:"required,hostname"`
	Port uint16 `yaml:"port"`
}

func (s TCPSocket) IsEmpty() bool {
	return s.Host == ""
}

func (s TCPSocket) Validate() error {
	if s.Host == "" {
		return fmt.Errorf("%w: host is empty", ErrValidation)
	}
	if s.Port == 0 {
		return fmt.Errorf("%w: port is empty", ErrValidation)
	}
	return nil
}
