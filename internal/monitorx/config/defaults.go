package config

var (
	defaultConfig Config = Config{
		Storage: Storage{
			DSN: "file://./config.yaml",
		},
	}
)
