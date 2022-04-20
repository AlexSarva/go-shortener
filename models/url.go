package models

import "time"

type Url struct {
	Id       string
	ShortUrl string
	RawUrl   string
	Created  time.Time
}
