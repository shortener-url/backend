package ports

import (
	"time"

	"github.com/shortener-url/backend/internal/core/domain"
)

type URLRepository interface {
	CreateURLShort(urlLarge, urlShort string) error
	GetURL(shortURL string) (*domain.URL, error)
	AddClick(idURL uint, ip, browser, referer string) error
	UpdateShortURL(shortURLOld, shortURLNew string) error
	ExistsShortURL(shorts []string) ([]*string, error)
	FindByParams(key string, startDate, endDate time.Time) ([]domain.URL, error)
}

type URLService interface {
	Create(url string, try, keyLen int) (*string, error)
	ExistsShortURL(shortURL string) (bool, error)
	UpdateShortURL(shortURLOld, shortURLNew string) (*string, error)
	Get(shortURL string, ip, browser, referer string) (*string, error)
	FindByParams(key string, startDate, endDate time.Time) ([]domain.URL, error)
}
