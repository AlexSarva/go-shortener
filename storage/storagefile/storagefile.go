package storagefile

import (
	"AlexSarva/go-shortener/models"
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
)

var ErrURLNotFound = errors.New("URL not found")
var ErrUserURLsNotFound = errors.New("no URLs found by such userID")

type fileStorage struct {
	file   string
	writer *bufio.Writer
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

func (f *fileStorage) InsertURL(id, rawURL, baseURL, userID string) error {
	URLData := models.URL{
		ID:       id,
		RawURL:   rawURL,
		ShortURL: baseURL + "/" + id,
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
			UserUrlInfo := &models.UserURL{
				ShortURL: URLInfo.ShortURL,
				RawURL:   URLInfo.RawURL,
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
