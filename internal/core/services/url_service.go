package services

import (
	"fmt"
	"time"

	"github.com/shortener-url/backend/internal/core/domain"
	"github.com/shortener-url/backend/internal/core/ports"
	"github.com/shortener-url/backend/utils"
	"gitlab.com/uchile1/helper/helperLog"
)

type urlService struct {
	repo ports.URLRepository
}

func NewURLService(repo ports.URLRepository) ports.URLService {
	return &urlService{repo: repo}
}

func (s *urlService) Create(url string, try, keyLen int) (*string, error) {
	var candidates []string
	for i := 0; i < try; i++ {
		candidates = append(candidates, utils.Base62Encode(keyLen))
	}
	availables, err := s.repo.ExistsShortURL(candidates)
	if err != nil {
		return nil, err
	}
	if len(availables) == 0 {
		return nil, fmt.Errorf("no available keys shortened")
	}
	err = s.repo.CreateURLShort(url, *availables[0])
	if err != nil {
		return nil, err
	}
	return availables[0], nil
}

func (s *urlService) ExistsShortURL(shortURL string) (bool, error) {
	urlsFound, err := s.repo.ExistsShortURL([]string{shortURL})
	if err != nil {
		return true, err
	}
	if len(urlsFound) == 1 {
		return false, nil
	}
	return true, nil
}

func (s *urlService) UpdateShortURL(shortURLOld, shortURLNew string) (*string, error) {
	exists, err := s.ExistsShortURL(shortURLNew)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("URL already exists")
	}
	err = s.repo.UpdateShortURL(shortURLOld, shortURLNew)
	if err != nil {
		return nil, err
	}
	return &shortURLNew, nil
}

func (s *urlService) Get(shortURL string, ip, browser, referer string) (*string, error) {
	url, err := s.repo.GetURL(shortURL)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, nil
	}

	go func() {
		err := s.repo.AddClick(url.ID, ip, browser, referer)
		if err != nil {
			helperLog.Logger.Error().Err(err).Msg("Error adding click in go routine")
		}
	}()
	return &url.OriginalURL, nil
}

func (s *urlService) FindByParams(key string, startDate, endDate time.Time) ([]domain.URL, error) {
	return s.repo.FindByParams(key, startDate, endDate)
}
