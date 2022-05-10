package utils

import (
	"bytes"
	"compress/gzip"
	"math/rand"
	"regexp"
	"time"
)

func ValidateURL(rawText string) bool {
	var re = regexp.MustCompile(`(\b(https?):\/\/)?[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+\.[-A-Za-z0-9+&@#\/%=~_|]+`)
	return re.Match([]byte(rawText))
}

func ValidateShortURL(rawText string) bool {
	var re = regexp.MustCompile(`http:\/\/localhost:8080\/[a-zA-Z]{5}`)
	return re.Match([]byte(rawText))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func ShortURLGenerator(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Compress сжимает слайс байт.
func GzipCompress(data []byte) []byte {
	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil
	}

	_, err = w.Write(data)
	if err != nil {
		return nil
	}

	err = w.Close()
	if err != nil {
		return nil
	}

	return b.Bytes()
}
