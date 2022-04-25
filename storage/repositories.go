package storage

import "AlexSarva/go-shortener/models"

type Repo interface {
	InsertURL(rawURL, shortURL string) error
	GetURL(shortURL string) (*models.URL, error)
}
