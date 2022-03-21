package conf

import (
	"fmt"
	"os"
)

type Config struct {
	ApiKey string
	SecretKey string
}

func LoadConfig() (*Config, error) {
	apiKey := os.Getenv("API_KEY")
	secret := os.Getenv("SECRET_KEY")
	if len(apiKey) == 0 || len(secret) == 0 {
		return nil, fmt.Errorf("required data is null")
	}

	return &Config{
		ApiKey:    apiKey,
		SecretKey: secret,
	}, nil
}