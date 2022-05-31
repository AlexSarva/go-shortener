package models

import "time"

type URL struct {
	ID       string    `json:"id" db:"id"`
	ShortURL string    `json:"short_url" db:"short_url"`
	RawURL   string    `json:"raw_url" db:"raw_url"`
	UserID   string    `json:"user_id" db:"user_id"`
	Created  time.Time `json:"created,omitempty" db:"created"`
}

type NewURL struct {
	URL string `json:"url"`
}

type ResultURL struct {
	Result string `json:"result"`
}

type UserURL struct {
	ShortURL string `json:"short_url" db:"short_url"`
	RawURL   string `json:"original_url" db:"raw_url"`
}

type AllUserURLs struct {
	URLList []UserURL
}
