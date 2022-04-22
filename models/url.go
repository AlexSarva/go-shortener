package models

import "time"

type URL struct {
	Id       string
	ShortURL string
	RawURL   string
	Created  time.Time
}
