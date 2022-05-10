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

type fileStorage struct {
	file string
}

func NewFileStorage(fileName string) (*fileStorage, error) {
	return &fileStorage{
		file: fileName,
	}, nil
}

func (f *fileStorage) InsertURL(id, rawURL, baseURL string) error {
	file, err := os.OpenFile(f.file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	closeErr := file.Close()
	if closeErr != nil {
		log.Fatal(closeErr)
	}

	writer := *bufio.NewWriter(file)

	URLData := models.URL{
		ID:       id,
		RawURL:   rawURL,
		ShortURL: "http://" + baseURL + "/" + id,
		Created:  time.Now(),
	}

	data, err := json.Marshal(&URLData)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := writer.Write(data); err != nil {
		return err
	}
	// добавляем перенос строки
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}
	// записываем буфер в файл
	return writer.Flush()
}

func (f *fileStorage) GetURL(id string) (*models.URL, error) {
	file, err := os.OpenFile(f.file, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	closeErr := file.Close()
	if closeErr != nil {
		log.Fatal(closeErr)
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
