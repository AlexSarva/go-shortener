package storage

import (
	"errors"
	"github.com/google/uuid"
	"go-shortener/models"
	"sync"
	"time"
)

var ErrUrlNotFound = errors.New("url not found")

type UrlLocalStorage struct {
	urlList map[string]*models.Url
	mutex   *sync.Mutex
}

func NewUrlLocalStorage() *UrlLocalStorage {
	return &UrlLocalStorage{
		urlList: make(map[string]*models.Url),
		mutex:   new(sync.Mutex),
	}
}

func (s *UrlLocalStorage) Insert(rawUrl, shortUrl string) error {
	id := uuid.New()
	urlData := &models.Url{
		Id:       id.String(),
		RawUrl:   rawUrl,
		ShortUrl: shortUrl,
		Created:  time.Now(),
	}
	s.mutex.Lock()
	s.urlList[shortUrl] = urlData
	s.mutex.Unlock()
	return nil
}

func (s *UrlLocalStorage) Get(shortUrl string) (*models.Url, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	urlInfo, ok := s.urlList[shortUrl]
	if ok == false {
		return &models.Url{}, ErrUrlNotFound
	}
	return urlInfo, nil
}
