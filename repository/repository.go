package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateURLShort(urlLarge, urlShort string) error {
	url := &URL{OriginalURL: urlLarge, ShortURL: urlShort}
	return r.DB.Create(url).Error
}

func (r *Repository) GetURL(shortURL string) (*URL, error) {
	var url URL
	err := r.DB.First(&url, &URL{ShortURL: shortURL}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &url, nil
}

func (r *Repository) AddClick(idURL uint, ip, browser, referer string) error {
	click := Click{URLID: idURL, IP: ip, Browser: browser, Referer: referer}
	return r.DB.Create(&click).Error
}

func (r *Repository) ExistsURLShort(shorts []string) (availables []*string, err error) {
	var urls []URL
	err = r.DB.Find(&urls).Where("short_url in ?", shorts).Error
	if err != nil {
		return nil, err
	}
	for _, candidate := range shorts {
		exists := false
		for _, url := range urls {
			if url.ShortURL == candidate {
				exists = true
				break
			}
		}
		if !exists {
			availables = append(availables, &candidate)
		}
	}
	return availables, nil
}

func (r *Repository) FindByParams(key string, startDate, endDate time.Time) ([]URL, error) {
	betweenDates := r.DB.Where("1=1")
	if !startDate.IsZero() && !endDate.IsZero() {
		betweenDates = r.DB.Where("click.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query := r.DB.Model(&URL{}).Distinct().Preload("Clicks", betweenDates)
	if key != "" {
		query = query.Where("short_url =?", key)
	}
	var c Click
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Joins(fmt.Sprintf("JOIN %s ON click.url_id = url.id", c.TableName())).
			Where(betweenDates)
	}
	var urls []URL

	err := query.Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}
