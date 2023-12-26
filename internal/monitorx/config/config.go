package config

type Config struct {
	Storage  Storage  `yaml:"storage"`
	EventBus EventBus `yaml:"eventBus"`
}

type Storage struct {
	DSN string `yaml:"dsn"`
}

type EventBus struct {
	DSN string `yaml:"dsn"`
}
