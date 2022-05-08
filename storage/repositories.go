package storage

import "AlexSarva/go-shortener/models"

type Repo interface {
	InsertURL(rawURL, shortURL, baseURL string) error
	GetURL(shortURL string) (*models.URL, error)
}
