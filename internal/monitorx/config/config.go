package config

type Config struct {
	Storage  Storage  `yaml:"storage"`
	EventBus EventBus `yaml:"eventBus"`
}

type Storage struct {
	DSN string `yaml:"dsn" envconfig:"STORAGE__DSN" validate:"required"`
}

type EventBus struct {
	DSN string `yaml:"dsn" envconfig:"EVENTBUS__DSN" validate:"required"`
}
