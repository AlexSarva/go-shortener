package checker

import "regexp"

func CheckUrl(rawText string) bool {
	var re = regexp.MustCompile(`(?m)https?:\/\/(www\.)?[-А-Яа-яa-zA-Z0-9@:%._\+\/?&~#=]+`)
	return re.Match([]byte(rawText))
}
