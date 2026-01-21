package services

import (
	"errors"
	"testing"
	"time"

	"github.com/shortener-url/backend/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockURLRepository is a mock implementation of ports.URLRepository
type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) CreateURLShort(urlLarge, urlShort string) error {
	args := m.Called(urlLarge, urlShort)
	return args.Error(0)
}

func (m *MockURLRepository) GetURL(shortURL string) (*domain.URL, error) {
	args := m.Called(shortURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.URL), args.Error(1)
}

func (m *MockURLRepository) AddClick(idURL uint, ip, browser, referer string) error {
	args := m.Called(idURL, ip, browser, referer)
	return args.Error(0)
}

func (m *MockURLRepository) UpdateShortURL(shortURLOld, shortURLNew string) error {
	args := m.Called(shortURLOld, shortURLNew)
	return args.Error(0)
}

func (m *MockURLRepository) ExistsShortURL(shorts []string) ([]*string, error) {
	args := m.Called(shorts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*string), args.Error(1)
}

func (m *MockURLRepository) FindByParams(key string, startDate, endDate time.Time) ([]domain.URL, error) {
	args := m.Called(key, startDate, endDate)
	return args.Get(0).([]domain.URL), args.Error(1)
}

func TestURLService_Create(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	t.Run("success", func(t *testing.T) {
		url := "https://example.com"
		mockRepo.On("ExistsShortURL", mock.Anything).Return([]*string{ptr("short")}, nil).Once()
		mockRepo.On("CreateURLShort", url, "short").Return(nil).Once()

		result, err := service.Create(url, 1, 5)

		assert.NoError(t, err)
		assert.Equal(t, "short", *result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("no available keys", func(t *testing.T) {
		url := "https://example.com"
		mockRepo.On("ExistsShortURL", mock.Anything).Return([]*string{}, nil).Once()

		result, err := service.Create(url, 1, 5)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "no available keys shortened", err.Error())
	})

	t.Run("repository error", func(t *testing.T) {
		url := "https://example.com"
		mockRepo.On("ExistsShortURL", mock.Anything).Return(nil, errors.New("db error")).Once()

		result, err := service.Create(url, 1, 5)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestURLService_ExistsShortURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	t.Run("exists", func(t *testing.T) {
		mockRepo.On("ExistsShortURL", []string{"exist"}).Return([]*string{}, nil).Once()

		exists, err := service.ExistsShortURL("exist")

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("does not exist", func(t *testing.T) {
		mockRepo.On("ExistsShortURL", []string{"notexist"}).Return([]*string{ptr("notexist")}, nil).Once()

		exists, err := service.ExistsShortURL("notexist")

		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestURLService_UpdateShortURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("ExistsShortURL", []string{"new"}).Return([]*string{ptr("new")}, nil).Once()
		mockRepo.On("UpdateShortURL", "old", "new").Return(nil).Once()

		result, err := service.UpdateShortURL("old", "new")

		assert.NoError(t, err)
		assert.Equal(t, "new", *result)
	})

	t.Run("already exists", func(t *testing.T) {
		mockRepo.On("ExistsShortURL", []string{"exists"}).Return([]*string{}, nil).Once()

		result, err := service.UpdateShortURL("old", "exists")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "URL already exists", err.Error())
	})
}

func TestURLService_Get(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	t.Run("success", func(t *testing.T) {
		short := "short"
		urlDomain := &domain.URL{ID: 1, OriginalURL: "https://example.com"}
		mockRepo.On("GetURL", short).Return(urlDomain, nil).Once()
		mockRepo.On("AddClick", uint(1), "127.0.0.1", "browser", "referer").Return(nil).Once()

		result, err := service.Get(short, "127.0.0.1", "browser", "referer")

		// wait a bit for go routine
		time.Sleep(100 * time.Millisecond)

		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", *result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetURL", "notfound").Return(nil, nil).Once()

		result, err := service.Get("notfound", "", "", "")

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func ptr(s string) *string {
	return &s
}
