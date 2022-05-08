package storage

import "AlexSarva/go-shortener/models"

type Repo interface {
	InsertURL(id, rawURL, baseURL string) error
	GetURL(id string) (*models.URL, error)
}
