package config

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram Telegram  `yaml:"telegram"`
	Services []Service `yaml:"services"`
	Storage  Storage   `yaml:"storage"`
}
type Telegram struct {
	Token      string `yaml:"token" envconfig:"TELEGRAM__TOKEN" validate:"required"`
	ChatID     int64  `yaml:"chatId" envconfig:"TELEGRAM__CHAT_ID"`
	WebhookURL string `yaml:"webhookUrl" envconfig:"TELEGRAM__WEBHOOK_URL" validate:"required"`
	Debug      bool   `yaml:"debug" envconfig:"TELEGRAM__DEBUG"`
}
type HTTPHeader struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
type HTTPGet struct {
	TCPSocket   `yaml:",inline"`
	Scheme      string       `yaml:"scheme"`
	Path        string       `yaml:"path"`
	HTTPHeaders []HTTPHeader `yaml:"httpHeaders"`
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

type Redis struct {
	Host     string `yaml:"host" envconfig:"REDIS__HOST"`
	Port     int    `yaml:"port" envconfig:"REDIS__PORT"`
	Password string `yaml:"password" envconfig:"REDIS__PASSWORD"`
	DB       int    `yaml:"db" envconfig:"REDIS__DB"`
}
type Storage struct {
	Redis Redis `yaml:"redis"`
}

var instance Config
var once = sync.Once{}

func loadConfig() Config {
	if err := godotenv.Load(".env"); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			errorLog.Println(err)
		}
	}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yml"
	}

	config := Config{}

	if err := fromYaml(path, &config); err != nil {
		errorLog.Printf("couldn'n load config from %s: %s\r\n", path, err.Error())
	}

	if err := fromEnv(&config); err != nil {
		errorLog.Printf("couldn'n load config from env: %s\r\n", err.Error())
	}

	return config
}

func fromYaml(path string, config *Config) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

func fromEnv(config *Config) error {
	return envconfig.Process("", config)
}

func GetConfig() Config {
	once.Do(func() {
		instance = loadConfig()
	})

	return instance
}
