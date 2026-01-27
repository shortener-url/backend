package repositories

import (
	"fmt"
	"time"

	"github.com/shortener-url/backend/internal/core/domain"
	"github.com/shortener-url/backend/internal/core/ports"
	"gorm.io/gorm"
)

type PostgresURL struct {
	ID          uint `gorm:"primary_key;autoIncrement"`
	OriginalURL string
	ShortURL    string
	Clicks      []PostgresClick `gorm:"foreignKey:URLID;references:ID"`
	CreatedAt   time.Time       `gorm:"<-:create"`
	UpdatedAt   time.Time
}

func (PostgresURL) TableName() string {
	return "shortener.url"
}

type PostgresClick struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	URLID     uint `gorm:"index"`
	IP        string
	Browser   string
	Referer   string
	CreatedAt time.Time `gorm:"<-:create"`
}

func (PostgresClick) TableName() string {
	return "shortener.click"
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) ports.URLRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateURLShort(urlLarge, urlShort string) error {
	url := &PostgresURL{OriginalURL: urlLarge, ShortURL: urlShort}
	return r.db.Create(url).Error
}

func (r *postgresRepository) GetURL(shortURL string) (*domain.URL, error) {
	var pgURL PostgresURL
	err := r.db.First(&pgURL, &PostgresURL{ShortURL: shortURL}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &domain.URL{
		ID:          pgURL.ID,
		OriginalURL: pgURL.OriginalURL,
		ShortURL:    pgURL.ShortURL,
		CreatedAt:   pgURL.CreatedAt,
		UpdatedAt:   pgURL.UpdatedAt,
	}, nil
}

func (r *postgresRepository) AddClick(idURL uint, ip, browser, referer string) error {
	click := PostgresClick{URLID: idURL, IP: ip, Browser: browser, Referer: referer}
	return r.db.Create(&click).Error
}

func (r *postgresRepository) UpdateShortURL(shortURLOld, shortURLNew string) error {
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	tx := r.db.Model(&PostgresURL{}).Where("short_url = ? AND created_at > ?", shortURLOld, oneHourAgo).Update("short_url", shortURLNew)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("no short URL found to update OR short URL was not created within the last hour")
	}
	return nil
}

func (r *postgresRepository) ExistsShortURL(shorts []string) (availables []*string, err error) {
	var urls []PostgresURL
	err = r.db.Where("short_url in ?", shorts).Find(&urls).Error
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
			c := candidate
			availables = append(availables, &c)
		}
	}
	return availables, nil
}

func (r *postgresRepository) FindByParams(key string, startDate, endDate time.Time) ([]domain.URL, error) {
	query := r.db.Model(&PostgresURL{}).Distinct().Preload("Clicks", func(db *gorm.DB) *gorm.DB {
		if !startDate.IsZero() && !endDate.IsZero() {
			db = db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
		return db.Order("created_at ASC")
	})
	if key != "" {
		query = query.Where("short_url =?", key)
	}

	var pgURLs []PostgresURL
	err := query.Find(&pgURLs).Error
	if err != nil {
		return nil, err
	}

	domainURLs := make([]domain.URL, len(pgURLs))
	for i, pgURL := range pgURLs {
		clicks := make([]domain.Click, len(pgURL.Clicks))
		for j, pgClick := range pgURL.Clicks {
			clicks[j] = domain.Click{
				ID:        pgClick.ID,
				URLID:     pgClick.URLID,
				IP:        pgClick.IP,
				Browser:   pgClick.Browser,
				Referer:   pgClick.Referer,
				CreatedAt: pgClick.CreatedAt,
			}
		}
		domainURLs[i] = domain.URL{
			ID:          pgURL.ID,
			OriginalURL: pgURL.OriginalURL,
			ShortURL:    pgURL.ShortURL,
			Clicks:      clicks,
			CreatedAt:   pgURL.CreatedAt,
			UpdatedAt:   pgURL.UpdatedAt,
		}
	}
	return domainURLs, nil
}
