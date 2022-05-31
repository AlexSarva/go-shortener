package storage

import (
	"AlexSarva/go-shortener/models"
)

type Repo interface {
	InsertURL(id, rawURL, baseURL, userID string) error
	GetURL(id string) (*models.URL, error)
	GetUserURLs(userID string) ([]models.UserURL, error)
	Ping() bool
}
