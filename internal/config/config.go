package config

import (
	"fmt"
	"math/rand"
	"strings"
)

type Config struct {
	Telegram Telegram  `yaml:"telegram"`
	Services []Service `yaml:"services"`
}
type Telegram struct {
	Token      string           `yaml:"token" envconfig:"TELEGRAM__TOKEN" validate:"required"`
	ChatID     int64            `yaml:"chatId" envconfig:"TELEGRAM__CHAT_ID"`
	WebhookURL string           `yaml:"webhookUrl" envconfig:"TELEGRAM__WEBHOOK_URL" validate:"required"`
	Debug      bool             `yaml:"debug" envconfig:"TELEGRAM__DEBUG"`
	Messages   TelegramMessages `yaml:"messages"`
}
type TelegramMessages map[string]string
type HTTPHeader struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type HTTPHeaders []HTTPHeader

func (h HTTPHeaders) ToMap() map[string][]string {
	m := make(map[string][]string, len(h))

	for _, v := range h {
		m[v.Name] = append(m[v.Name], v.Value)
	}

	return m
}

type HTTPGet struct {
	TCPSocket   `yaml:",inline"`
	Scheme      string      `yaml:"scheme"`
	Path        string      `yaml:"path"`
	HTTPHeaders HTTPHeaders `yaml:"httpHeaders"`
}

func (s HTTPGet) IsEmpty() bool {
	return s.Host == ""
}

func (h HTTPGet) ApplyDefaultsAndValidate() (HTTPGet, error) {
	if h.IsEmpty() {
		return h, nil
	}

	if h.Scheme == "" {
		if h.Port == 80 {
			h.Scheme = "http"
		}
		if h.Port == 443 {
			h.Scheme = "https"
		}
	}

	if h.Port == 0 {
		if h.Scheme == "http" {
			h.Port = 80
		}
		if h.Scheme == "https" {
			h.Port = 443
		}
	}

	if h.Scheme != "http" && h.Scheme != "https" {
		return h, fmt.Errorf("invalid scheme %s", h.Scheme)
	}

	if !strings.HasPrefix(h.Path, "/") {
		h.Path = "/" + h.Path
	}

	return h, nil
}

type TCPSocket struct {
	Host string `yaml:"host" validate:"required,hostname"`
	Port uint16 `yaml:"port"`
}

func (s TCPSocket) IsEmpty() bool {
	return s.Host == ""
}

type Service struct {
	Id                  string    `yaml:"id"`
	Name                string    `yaml:"name"`
	InitialDelaySeconds int       `yaml:"initialDelaySeconds"` // пауза перед первым опросом в секундах, по умолчанию: 0; если меньше 0, то используется случайное значение между 0 и `periodSeconds`
	PeriodSeconds       int       `yaml:"periodSeconds"`       // период опроса в секундах, по умолчанию: 10
	TimeoutSeconds      int       `yaml:"timeoutSeconds"`      // время ожидания ответа в секундах, по кмолчанию: 1
	SuccessThreshold    int       `yaml:"successThreshold"`    // количество последовательных успешных соединений для перехода в состояние "в сети", по умолчанию: 1
	FailureThreshold    int       `yaml:"failureThreshold"`    // количество последовательных ошибок соединения для перехода в состояние "не в сети", по умолчанию: 3
	HTTPGet             HTTPGet   `yaml:"httpGet,omitempty"`
	TCPSocket           TCPSocket `yaml:"tcpSocket,omitempty"`
}

func (s Service) ApplyDefaultsAndValidate() (svc Service, err error) {
	if s.PeriodSeconds < 0 {
		return s, fmt.Errorf("periodSeconds must be gte 0")
	}
	if s.SuccessThreshold < 0 {
		return s, fmt.Errorf("successThreshold must be gte 0")
	}
	if s.FailureThreshold < 0 {
		return s, fmt.Errorf("failureThreshold must be gte 0")
	}
	if s.HTTPGet.IsEmpty() && s.TCPSocket.IsEmpty() {
		return s, fmt.Errorf("one of httpGet or tcpSocket must be filled")
	}

	if s.HTTPGet, err = s.HTTPGet.ApplyDefaultsAndValidate(); err != nil {
		return s, err
	}

	if s.PeriodSeconds == 0 {
		s.PeriodSeconds = 10
	}

	if s.InitialDelaySeconds < 0 {
		s.InitialDelaySeconds = rand.Intn(s.PeriodSeconds + 1)
	}

	if s.SuccessThreshold == 0 {
		s.SuccessThreshold = 1
	}

	if s.FailureThreshold == 0 {
		s.FailureThreshold = 3
	}

	return s, nil
}

type Monitor struct {
	Services []Service `yaml:"services"`
}
