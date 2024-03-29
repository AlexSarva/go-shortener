package storagefile

import (
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/utils"
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
)

var ErrURLNotFound = errors.New("URL not found")
var ErrUserURLsNotFound = errors.New("no URLs found by such userID")
var ErrEmptyData = errors.New("no data in file")

type fileStorage struct {
	writer *bufio.Writer
	file   string
}

func NewFileStorage(fileName string) (*fileStorage, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return &fileStorage{
		file:   fileName,
		writer: bufio.NewWriter(file),
	}, nil
}

func (f *fileStorage) Ping() bool {
	return true
}

func (f *fileStorage) InsertURL(id, rawURL, shortURL, userID string) error {
	URLData := models.URL{
		ID:       id,
		RawURL:   rawURL,
		ShortURL: shortURL,
		Created:  time.Now(),
		UserID:   userID,
	}

	data, err := json.Marshal(&URLData)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := f.writer.Write(data); err != nil {
		return err
	}
	// добавляем перенос строки
	if err := f.writer.WriteByte('\n'); err != nil {
		return err
	}
	// записываем буфер в файл

	return f.writer.Flush()
}

func (f *fileStorage) InsertMany(bathURL []models.URL) error {
	for _, curURL := range bathURL {
		data, err := json.Marshal(&curURL)
		if err != nil {
			return err
		}
		if _, err := f.writer.Write(data); err != nil {
			return err
		}
		if err := f.writer.WriteByte('\n'); err != nil {
			return err
		}
	}
	return f.writer.Flush()
}

func (f *fileStorage) GetURL(id string) (*models.URL, error) {
	file, err := os.OpenFile(f.file, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	scanner := *bufio.NewScanner(file)

	for scanner.Scan() {
		var URLInfo models.URL
		// читаем данные из scanner
		data := scanner.Bytes()
		if err := json.Unmarshal(data, &URLInfo); err != nil {
			panic(err)
		}
		if URLInfo.ID == id {
			log.Printf("%+v\n", URLInfo)
			return &URLInfo, nil
		}
	}

	return &models.URL{}, ErrURLNotFound
}

func (f *fileStorage) GetURLByRaw(rawURL string) (*models.URL, error) {
	file, err := os.OpenFile(f.file, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	scanner := *bufio.NewScanner(file)

	for scanner.Scan() {
		var URLInfo models.URL
		// читаем данные из scanner
		data := scanner.Bytes()
		if err := json.Unmarshal(data, &URLInfo); err != nil {
			panic(err)
		}
		if URLInfo.RawURL == rawURL {
			return &URLInfo, nil
		}
	}
	return &models.URL{}, ErrURLNotFound
}

func (f *fileStorage) GetUserURLs(userID string) ([]models.UserURL, error) {
	file, err := os.OpenFile(f.file, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	scanner := *bufio.NewScanner(file)
	var URLList []models.UserURL
	for scanner.Scan() {
		var URLInfo models.URL
		// читаем данные из scanner
		data := scanner.Bytes()
		if err := json.Unmarshal(data, &URLInfo); err != nil {
			panic(err)
		}
		if URLInfo.UserID == userID {
			UserURLInfo := &models.UserURL{
				ShortURL: URLInfo.ShortURL,
				RawURL:   URLInfo.RawURL,
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

func (f *fileStorage) Delete(userID string, shortURLs []string) error {
	return nil
}

func (f *fileStorage) GetStat() (*models.SystemStat, error) {
	var stat models.SystemStat
	file, err := os.OpenFile(f.file, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	scanner := *bufio.NewScanner(file)
	var usersList []string
	for scanner.Scan() {
		var URLInfo models.URL
		// читаем данные из scanner
		data := scanner.Bytes()
		if err := json.Unmarshal(data, &URLInfo); err != nil {
			panic(err)
		}
		usersList = append(usersList, URLInfo.UserID)
	}
	stat.URLsCnt = len(usersList)
	stat.UsersCnt = len(utils.UniqueElements(usersList))

	if (models.SystemStat{}) == stat {
		return nil, ErrEmptyData
	}

	return &stat, nil
}
