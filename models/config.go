package models

// Config  start parameters for lunch the service
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage   string `env:"FILE_STORAGE_PATH"`
	Database      string `env:"DATABASE_DSN"`
}
