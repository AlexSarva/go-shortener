package models

import "time"

type URL struct {
	ID       string    `json:"id"`
	ShortURL string    `json:"short_url"`
	RawURL   string    `json:"raw_url"`
	Created  time.Time `json:"created,omitempty"`
}

type NewURL struct {
	URL string `json:"url"`
}

type ResultURL struct {
	Result string `json:"result"`
}
