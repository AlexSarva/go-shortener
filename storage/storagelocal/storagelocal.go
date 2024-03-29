package storagelocal

import (
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/utils"
	"errors"
	"log"
	"sync"
	"time"
)

var ErrURLNotFound = errors.New("URL not found")
var ErrUserURLsNotFound = errors.New("no URLs found by such userID")
var ErrEmptyData = errors.New("no data in DB")

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

func (s *URLLocalStorage) Ping() bool {
	return true
}

func (s *URLLocalStorage) InsertURL(id, rawURL, shortURL, userID string) error {
	URLData := &models.URL{
		ID:       id,
		RawURL:   rawURL,
		ShortURL: shortURL,
		Created:  time.Now(),
		UserID:   userID,
	}
	s.mutex.Lock()
	s.URLList[id] = URLData
	s.mutex.Unlock()
	log.Println("123")
	return nil
}

func (s *URLLocalStorage) InsertMany(bathURL []models.URL) error {
	s.mutex.Lock()
	for _, curURL := range bathURL {
		s.URLList[curURL.ID] = &curURL
	}
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

func (s *URLLocalStorage) GetURLByRaw(rawURL string) (*models.URL, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, value := range s.URLList {
		if value.RawURL == rawURL {
			return value, nil
		}
	}
	return &models.URL{}, ErrURLNotFound
}

func (s *URLLocalStorage) GetUserURLs(userID string) ([]models.UserURL, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var URLList []models.UserURL
	for _, urlInfo := range s.URLList {
		if urlInfo.UserID == userID {
			UserURLInfo := &models.UserURL{
				ShortURL: urlInfo.ShortURL,
				RawURL:   urlInfo.RawURL,
			}
			URLList = append(URLList, *UserURLInfo)
		}
	}
	if len(URLList) > 0 {
		return URLList, nil
	} else {
		return URLList, ErrUserURLsNotFound
	}
}

func (s *URLLocalStorage) Delete(userID string, shortURLs []string) error {
	return nil
}

func (s *URLLocalStorage) GetStat() (*models.SystemStat, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var stat models.SystemStat
	stat.URLsCnt = len(s.URLList)

	var userList []string
	for _, urlInfo := range s.URLList {
		userList = append(userList, urlInfo.UserID)
	}
	stat.UsersCnt = len(utils.UniqueElements(userList))

	if (models.SystemStat{}) == stat {
		return nil, ErrEmptyData
	}

	return &stat, nil
}
