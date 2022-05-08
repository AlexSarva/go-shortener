package storagelocal

import (
	"AlexSarva/go-shortener/models"

	"errors"
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

func (s *URLLocalStorage) InsertURL(id, rawURL, baseURL string) error {
	URLData := &models.URL{
		ID:       id,
		RawURL:   rawURL,
		ShortURL: "http://" + baseURL + "/" + id,
		Created:  time.Now(),
	}
	s.mutex.Lock()
	s.URLList[id] = URLData
	s.mutex.Unlock()
	return nil
}

func (s *URLLocalStorage) GetURL(id string) (*models.URL, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	URLInfo, ok := s.URLList[id]
	if !ok {
		return &models.URL{}, ErrURLNotFound
	}
	return URLInfo, nil
}
