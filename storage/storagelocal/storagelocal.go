package storagelocal

import (
	"AlexSarva/go-shortener/models"

	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
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

func (s *URLLocalStorage) InsertURL(rawURL, shortURL, baseURL string) error {
	id := uuid.New()
	URLData := &models.URL{
		ID:       id.String(),
		RawURL:   rawURL,
		ShortURL: "http://" + baseURL + "/" + shortURL,
		Created:  time.Now(),
	}
	s.mutex.Lock()
	s.URLList[shortURL] = URLData
	s.mutex.Unlock()
	return nil
}

func (s *URLLocalStorage) GetURL(shortURL string) (*models.URL, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	URLInfo, ok := s.URLList[shortURL]
	if !ok {
		return &models.URL{}, ErrURLNotFound
	}
	return URLInfo, nil
}
