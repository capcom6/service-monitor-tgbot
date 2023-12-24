package storage

import (
	"errors"
	"math/rand"
)

type Config struct {
	DSN string
}

type Service struct {
	Id                  string    `yaml:"id" json:"id" validate:"required"`
	Name                string    `yaml:"name" json:"name" validate:"required"`
	InitialDelaySeconds int32     `yaml:"initialDelaySeconds" json:"initialDelaySeconds"`           // пауза перед первым опросом в секундах, по умолчанию: 0; если меньше 0, то используется случайное значение между 0 и `periodSeconds`
	PeriodSeconds       uint16    `yaml:"periodSeconds" json:"periodSeconds" validate:"gt=0"`       // период опроса в секундах, по умолчанию: 10
	TimeoutSeconds      uint16    `yaml:"timeoutSeconds" json:"timeoutSeconds" validate:"gt=0"`     // время ожидания ответа в секундах, по кмолчанию: 1
	SuccessThreshold    uint8     `yaml:"successThreshold" json:"successThreshold" validate:"gt=0"` // количество последовательных успешных соединений для перехода в состояние "в сети", по умолчанию: 1
	FailureThreshold    uint8     `yaml:"failureThreshold" json:"failureThreshold" validate:"gt=0"` // количество последовательных ошибок соединения для перехода в состояние "не в сети", по умолчанию: 3
	HTTPGet             HTTPGet   `yaml:"httpGet,omitempty" json:"httpGet,omitempty" validate:"required_without=TCPSocket"`
	TCPSocket           TCPSocket `yaml:"tcpSocket,omitempty" json:"tcpSocket,omitempty" validate:"required_without=HTTPGet"`
}

func (s *Service) Validate() error {
	if s.Id == "" {
		s.Id = s.Name
	}
	if s.Name == "" {
		return errors.New("name is empty")
	}
	if s.PeriodSeconds == 0 {
		s.PeriodSeconds = 10
	}
	if s.TimeoutSeconds == 0 {
		s.TimeoutSeconds = 1
	}
	if s.InitialDelaySeconds < 0 {
		s.InitialDelaySeconds = int32(rand.Intn(int(s.PeriodSeconds) + 1))
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

	return errors.New("one of httpGet or tcpSocket must be filled")
}

type HTTPGet struct {
	TCPSocket   `yaml:",inline" json:",inline"`
	Scheme      string      `yaml:"scheme" json:"scheme"`
	Path        string      `yaml:"path" json:"path"`
	HTTPHeaders HTTPHeaders `yaml:"httpHeaders" json:"httpHeaders"`
}

func (s *HTTPGet) Validate() error {
	if s.Scheme == "" {
		s.Scheme = "http"
	}
	if s.Path == "" {
		s.Path = "/"
	}

	if s.TCPSocket.Port == 0 {
		if s.Scheme == "http" {
			s.TCPSocket.Port = 80
		} else if s.Scheme == "https" {
			s.TCPSocket.Port = 443
		} else {
			return errors.New("port is empty")
		}
	}

	return s.TCPSocket.Validate()
}

type HTTPHeaders []HTTPHeader

type HTTPHeader struct {
	Name  string `yaml:"name" json:"name" validate:"required"`
	Value string `yaml:"value" json:"value" validate:"required"`
}

type TCPSocket struct {
	Host string `yaml:"host" json:"host" validate:"required,hostname"`
	Port uint16 `yaml:"port" json:"port" validate:"gt=0"`
}

func (s TCPSocket) IsEmpty() bool {
	return s.Host == ""
}

func (s *TCPSocket) Validate() error {
	if s.Host == "" {
		return errors.New("host is empty")
	}
	if s.Port == 0 {
		return errors.New("port is empty")
	}
	return nil
}
