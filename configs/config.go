package configs

import "os"

type Config struct {
	BrevoAPIKey string
	FromEmail   string
}

func GetConfig() *Config {
	return &Config{
		BrevoAPIKey: os.Getenv("BREVO_API_KEY"),
		FromEmail:   os.Getenv("FROM_MAIL"),
	}
}
