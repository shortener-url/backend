package repository

import (
	"time"
)

type URL struct {
	ID          uint      `json:"url_ID" gorm:"primary_key;autoIncrement"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	Clicks      []Click   `json:"clicks" gorm:"foreignKey:URLID;references:ID"`
	CreatedAt   time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (URL) TableName() string {
	return "shortener.url"
}

type Click struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	URLID     uint      `json:"url_ID"  gorm:"index"`
	IP        string    `json:"ip"`
	Browser   string    `json:"browser"`
	Referer   string    `json:"referer"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create"`
}

func (Click) TableName() string {
	return "shortener.click"
}
