package models

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"localhost:8080"`
	FileStorage   string `env:"FILE_STORAGE_PATH"`
}
