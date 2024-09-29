package config

import (
	"musPlayer/models"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database   models.DatabaseConfig
	API        models.APIConfig
	Logging    models.LoggingConfig
	App        models.AppConfig
	SpotifyApi models.SpotifyConfig
}

func MustLoad() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Database: models.DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		API: models.APIConfig{
			BaseURL: os.Getenv("API_BASE_URL"),
			Key:     os.Getenv("API_KEY"),
		},
		Logging: models.LoggingConfig{
			Level: os.Getenv("LOG_LEVEL"),
		},
		App: models.AppConfig{
			Port: os.Getenv("APP_PORT"),
		},
		SpotifyApi: models.SpotifyConfig{
			ID:          os.Getenv("CLIENT_ID"),
			Secret:      os.Getenv("CLIENT_SECRET"),
			RedirectURI: os.Getenv("REDIRECT_URI"),
			AuthURL:     os.Getenv("AUTH_URL"),
			TokenURL:    os.Getenv("TOKEN_URL"),
			Scope:       os.Getenv("SCOPE"),
		},
	}

	return &cfg, nil
}
