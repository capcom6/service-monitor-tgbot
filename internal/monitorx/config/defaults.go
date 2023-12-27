package config

var (
	defaultConfig Config = Config{
		Storage: Storage{DSN: "file://./config.yaml"},
		EventBus: EventBus{
			DSN: "redis://localhost:6379/0?channel=events",
		},
	}
)
