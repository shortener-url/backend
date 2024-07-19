package model

type ShortenRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
}
