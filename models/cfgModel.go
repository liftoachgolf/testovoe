package models

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type APIConfig struct {
	BaseURL string
	Key     string
}

type LoggingConfig struct {
	Level string
}

type AppConfig struct {
	Port string
}

type GeniusConfig struct {
	ID          string
	Secret      string
	Token       string
	RedirectURI string
	AuthURL     string
	TokenURL    string
	Scope       string
}
