package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shortener-url/backend/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockURLService is a mock implementation of ports.URLService
type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) Create(url string, try, keyLen int) (*string, error) {
	args := m.Called(url, try, keyLen)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockURLService) ExistsShortURL(shortURL string) (bool, error) {
	args := m.Called(shortURL)
	return args.Bool(0), args.Error(1)
}

func (m *MockURLService) UpdateShortURL(shortURLOld, shortURLNew string) (*string, error) {
	args := m.Called(shortURLOld, shortURLNew)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockURLService) Get(shortURL string, ip, browser, referer string) (*string, error) {
	args := m.Called(shortURL, ip, browser, referer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockURLService) FindByParams(key string, startDate, endDate time.Time) ([]domain.URL, error) {
	args := m.Called(key, startDate, endDate)
	return args.Get(0).([]domain.URL), args.Error(1)
}

func TestHTTPHandler_Shorten(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockURLService)
	handler := NewHTTPHandler(mockSvc)
	r := gin.Default()
	r.POST("/shorten", handler.Shorten)

	t.Run("success", func(t *testing.T) {
		reqBody := domain.ShortenRequest{OriginalURL: "https://example.com"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		short := "abcde"
		mockSvc.On("Create", "https://example.com", mock.Anything, mock.Anything).Return(&short, nil).Once()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid url", func(t *testing.T) {
		reqBody := domain.ShortenRequest{OriginalURL: "invalid"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHTTPHandler_GetUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockURLService)
	handler := NewHTTPHandler(mockSvc)
	r := gin.Default()
	r.GET("/:key-short", handler.GetUrl)

	t.Run("success", func(t *testing.T) {
		short := "abcde"
		long := "https://example.com"
		mockSvc.On("Get", short, mock.Anything, mock.Anything, mock.Anything).Return(&long, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/"+short, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		short := "notfound"
		mockSvc.On("Get", short, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/"+short, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestHTTPHandler_ExistsShortURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockURLService)
	handler := NewHTTPHandler(mockSvc)
	r := gin.Default()
	r.GET("/:key-short/exists", handler.ExistsShortURL)

	t.Run("exists", func(t *testing.T) {
		short := "abcde"
		mockSvc.On("ExistsShortURL", short).Return(true, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/"+short+"/exists", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
