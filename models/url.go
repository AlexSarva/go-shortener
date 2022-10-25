package models

import "time"

// URL base struct for URLs in DB
type URL struct {
	Created  time.Time `json:"created,omitempty" db:"created"`
	ID       string    `json:"id" db:"id"`
	ShortURL string    `json:"short_url" db:"short_url"`
	RawURL   string    `json:"raw_url" db:"raw_url"`
	UserID   string    `json:"user_id" db:"user_id"`
	Deleted  int       `json:"deleted" db:"deleted"`
}

type RawBatchURL struct {
	CorrelationID string `json:"correlation_id" db:"id"`
	RawURL        string `json:"original_url" db:"raw_url"`
}

type ResultBatchURL struct {
	CorrelationID string `json:"correlation_id" db:"id"`
	ShortURL      string `json:"short_url" db:"short_url"`
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

type DeleteURL struct {
	UserID string
	URLs   []string
}

type SystemStat struct {
	URLsCnt  int `json:"urls" db:"urls_cnt"`
	UsersCnt int `json:"users" db:"users_cnt"`
}
