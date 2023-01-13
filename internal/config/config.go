package config

import (
	"errors"
	"io/fs"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram Telegram   `yaml:"telegram"`
	Services []Services `yaml:"services"`
	Storage  Storage    `yaml:"storage"`
}
type Telegram struct {
	Token      string `yaml:"token" envconfig:"TELEGRAM__TOKEN" validate:"required"`
	ChatID     int64  `yaml:"chatId" envconfig:"TELEGRAM__CHAT_ID"`
	WebhookURL string `yaml:"webhookUrl" envconfig:"TELEGRAM__WEBHOOK_URL" validate:"required"`
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

type TCPSocket struct {
	Host string `yaml:"host" validate:"required,hostname"`
	Port uint16 `yaml:"port"`
}

func (s TCPSocket) IsEmpty() bool {
	return s.Host == ""
}

type Services struct {
	Name                string    `yaml:"name"`
	InitialDelaySeconds int       `yaml:"initialDelaySeconds"`
	PeriodSeconds       int       `yaml:"periodSeconds"`
	TimeoutSeconds      int       `yaml:"timeoutSeconds"`
	SuccessThreshold    int       `yaml:"successThreshold"`
	FailureThreshold    int       `yaml:"failureThreshold"`
	HTTPGet             HTTPGet   `yaml:"httpGet,omitempty"`
	TCPSocket           TCPSocket `yaml:"tcpSocket,omitempty"`
}
type Redis struct {
	Host string `yaml:"host" envconfig:"REDIS__HOST"`
	Port int    `yaml:"port" envconfig:"REDIS__PORT"`
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
