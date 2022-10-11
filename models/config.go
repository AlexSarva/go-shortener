package models

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

// Config  start parameters for lunch the service
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	FileStorage   string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	Database      string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS   bool   `env:"ENABLE_HTTPS" json:"enable_https"`
}

// JSONConfig config file in json format
type JSONConfig struct {
	DSN string
}

func ReadJSONConfig(cfg *Config, JSONFilepath string) error {
	f, fErr := os.Open(JSONFilepath)
	if fErr != nil {
		log.Println("error opening file: err:", fErr)
	}
	defer f.Close()

	var unmarshalConfigErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&cfg)

	if err != nil {
		if errors.As(err, &unmarshalConfigErr) {
			log.Println("Wrong json format ")
			return err
		}
	}

	return nil
}
