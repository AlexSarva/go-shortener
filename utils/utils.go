package utils

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ValidateURL check original url by pattern
func ValidateURL(rawText string) bool {
	var re = regexp.MustCompile(`(\b(https?):\/\/)?[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+\.[-A-Za-z0-9+&@#\/%=~_|]+`)
	return re.Match([]byte(rawText))
}

// CreateShortURL create short url by concat base url of the service and generated id
func CreateShortURL(path, shortPath string) string {
	return fmt.Sprintf("%s/%s", path, shortPath)
}

// ValidateShortURL check short url by pattern
func ValidateShortURL(rawText, path string, n int) bool {
	pattern := fmt.Sprintf("%s/[a-zA-Z]{%d}", path, n)
	re := regexp.MustCompile(pattern)
	return re.Match([]byte(rawText))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// ShortURLGenerator generate id of the short url with the specified number of characters
func ShortURLGenerator(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func UniqueElements(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func ResolveIP(r *http.Request) (net.IP, error) {
	ipStr := r.Header.Get("X-Real-IP")
	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// если заголовок X-Real-IP пуст, пробуем X-Forwarded-For
		// этот заголовок содержит адреса отправителя и промежуточных прокси
		// в виде 203.0.113.195, 70.41.3.18, 150.172.238.178
		ips := r.Header.Get("X-Forwarded-For")
		// разделяем цепочку адресов
		ipStrs := strings.Split(ips, ",")
		// интересует только первый
		ipStr = ipStrs[0]
		// парсим
		ip = net.ParseIP(ipStr)
	}

	if ip == nil {
		addr := r.RemoteAddr
		ipStr2, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		ip = net.ParseIP(ipStr2)
		if ip == nil {
			return nil, fmt.Errorf("unexpected parse ip error")
		}
	}

	if ip == nil {
		return nil, fmt.Errorf("failed parse ip from http header")
	}

	return ip, nil

}
