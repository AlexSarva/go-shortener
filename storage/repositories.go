package storage

import (
	"AlexSarva/go-shortener/models"
	"errors"
)

var ERRDuplicatePK = errors.New("duplicate PK")

type Repo interface {
	InsertURL(id, rawURL, baseURL, userID string) error
	InsertMany(bathURL []models.URL) error
	GetURL(id string) (*models.URL, error)
	GetUserURLs(userID string) ([]models.UserURL, error)
	Ping() bool
	GetURLByRaw(rawURL string) (*models.URL, error)
}
