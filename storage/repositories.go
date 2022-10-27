package storage

import (
	"AlexSarva/go-shortener/models"
	"errors"
)

// ErrDuplicatePK expect when original url exists in DB
var ErrDuplicatePK = errors.New("duplicate PK")

// ErrNoValues error that occurs when no values selected from database
var ErrNoValues = errors.New("no values from select")

type Repo interface {
	// InsertURL add new link to DB
	InsertURL(id, rawURL, shortURL, userID string) error
	// InsertMany add several links to DB
	InsertMany(bathURL []models.URL) error
	// GetURL get original url from DB
	GetURL(id string) (*models.URL, error)
	// GetUserURLs get all user's urls from DB
	GetUserURLs(userID string) ([]models.UserURL, error)
	// Ping check connection to DB
	Ping() bool
	// GetURLByRaw get exist short url by original
	GetURLByRaw(rawURL string) (*models.URL, error)
	// Delete delete url from DB
	Delete(userID string, urls []string) error
	// GetStat get stata of users cnt ans urls cnt
	GetStat() (*models.SystemStat, error)
}
