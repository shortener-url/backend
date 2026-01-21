package domain

import "time"

type URL struct {
	ID          uint      `json:"url_ID"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	Clicks      []Click   `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Click struct {
	ID        uint      `json:"id"`
	URLID     uint      `json:"url_ID"`
	IP        string    `json:"ip"`
	Browser   string    `json:"browser"`
	Referer   string    `json:"referer"`
	CreatedAt time.Time `json:"created_at"`
}

type ShortenRequest struct {
	OriginalURL string `json:"original_url" example:"https://www.google.com"`
}

type URLUpdateRequest struct {
	ShortURLNew string `json:"short_url_new" example:"newkey"`
	ShortURLOld string `json:"short_url_old" example:"oldkey"`
}
