package models

import "time"

type URL struct {
	ID       string
	ShortURL string
	RawURL   string
	Created  time.Time
}
