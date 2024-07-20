package service

import (
	"fmt"
	"time"

	"github.com/shortener-url/backend/repository"
	"github.com/shortener-url/backend/utils"
	"gitlab.com/uchile1/helper/helperLog"
)

type Service struct {
	Repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Repository: repo}
}

func (service *Service) Create(url string, try, keyLen int) (*string, error) {
	var candidates []string
	for i := 0; i < int(try); i++ {
		candidates = append(candidates, utils.Base62Encode(keyLen))
	}
	availables, err := service.Repository.ExistsURLShort(candidates)
	if err != nil {
		return nil, err
	}
	if len(availables) == 0 {
		return nil, fmt.Errorf("no available keys shortened")
	}
	err = service.Repository.CreateURLShort(url, *availables[0])
	if err != nil {
		return nil, err
	}
	return availables[0], nil
}

func (service *Service) Get(shortURL string, ip, browser, referer string) (*string, error) {
	url, err := service.Repository.GetURL(shortURL)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, fmt.Errorf("url not found")
	}
	fmt.Println("antes de agregar click")
	go func() {
		err := service.Repository.AddClick(url.ID, ip, browser, referer)
		if err != nil {
			helperLog.Logger.Error().Err(err).Msg("Error adding click in go routine")
		}
	}()
	return &url.OriginalURL, nil
}

func (service *Service) FindByParams(key string, startDate, endDate time.Time) ([]repository.URL, error) {
	return service.Repository.FindByParams(key, startDate, endDate)
}
