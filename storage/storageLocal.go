package storage

import (
	"AlexSarva/go-shortener/models"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

var ErrURLNotFound = errors.New("URL not found")

type URLLocalStorage struct {
	URLList map[string]*models.URL
	mutex   *sync.Mutex
}

func NewURLLocalStorage() *URLLocalStorage {
	return &URLLocalStorage{
		URLList: make(map[string]*models.URL),
		mutex:   new(sync.Mutex),
	}
}

func (s *URLLocalStorage) Insert(rawURL, shortURL string) error {
	id := uuid.New()
	URLData := &models.URL{
		Id:       id.String(),
		RawURL:   rawURL,
		ShortURL: shortURL,
		Created:  time.Now(),
	}
	s.mutex.Lock()
	s.URLList[shortURL] = URLData
	s.mutex.Unlock()
	return nil
}

func (s *URLLocalStorage) Get(shortURL string) (*models.URL, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	URLInfo, ok := s.URLList[shortURL]
	if !ok {
		return &models.URL{}, ErrURLNotFound
	}
	return URLInfo, nil
}

func InitDB() *URLLocalStorage {
	db := *NewURLLocalStorage()
	return &db
}
