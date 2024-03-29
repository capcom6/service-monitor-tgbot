package config

var (
	defaultTelegramMessages TelegramMessages = TelegramMessages{
		"online":  "✅ {{.Name}} is *online*",
		"offline": "❌ {{.Name}} is *offline*: {{.Error}}",
	}

	defaultConfig Config = Config{
		Telegram: Telegram{
			Token:      "",
			ChatID:     0,
			WebhookURL: "",
			Debug:      false,
			Messages:   defaultTelegramMessages,
		},
		// Services: []Service{},
	}
)
