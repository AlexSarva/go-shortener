package storage

import "AlexSarva/go-shortener/models"

type Repo interface {
	InsertURL(id, rawURL, baseURL string, userID int) error
	GetURL(id string) (*models.URL, error)
	GetUserURLs(userID int) ([]models.UserURL, error)
}
