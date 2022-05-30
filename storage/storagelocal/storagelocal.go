package storagelocal

import (
	"AlexSarva/go-shortener/models"
	"errors"
	"sync"
	"time"
)

var ErrURLNotFound = errors.New("URL not found")
var ErrUserURLsNotFound = errors.New("no URLs found by such userID")

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

func (s *URLLocalStorage) InsertURL(id, rawURL, baseURL, userID string) error {
	URLData := &models.URL{
		ID:       id,
		RawURL:   rawURL,
		ShortURL: baseURL + "/" + id,
		Created:  time.Now(),
		UserID:   userID,
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

func (s *URLLocalStorage) GetUserURLs(userID string) ([]models.UserURL, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var URLList []models.UserURL
	for _, urlInfo := range s.URLList {
		if urlInfo.UserID == userID {
			UserUrlInfo := &models.UserURL{
				ShortURL: urlInfo.ShortURL,
				RawURL:   urlInfo.RawURL,
			}
			URLList = append(URLList, *UserUrlInfo)
		}
	}
	if len(URLList) > 0 {
		return URLList, nil
	} else {
		return URLList, ErrUserURLsNotFound
	}
}
