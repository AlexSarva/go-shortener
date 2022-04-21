package checker

import "regexp"

func CheckUrl(rawText string) bool {
	var re = regexp.MustCompile(`(\b(https?):\/\/)?[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+\.[-A-Za-z0-9+&@#\/%=~_|]+`)
	return re.Match([]byte(rawText))
}

func CheckShortUrl(rawText string) bool {
	var re = regexp.MustCompile(`http:\/\/localhost:8080\/[a-zA-Z]{5}`)
	return re.Match([]byte(rawText))
}
