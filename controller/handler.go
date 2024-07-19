package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shortener-url/backend/model"
	"github.com/shortener-url/backend/service"
	"gitlab.com/uchile1/helper/helperLog"
)

type Controller struct {
	Service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{Service: service}
}

func (controller *Controller) Shorten(c *gin.Context) {
	var req model.ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < 10; i++ {
		base62Encode()
	}

	c.Status(http.StatusNotImplemented)
}

func (controller *Controller) GetUrl(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (controller *Controller) Stats(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func base62Encode() string {
	timestamp := time.Now().UnixNano()

	const base = 62
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!()*+-/<=>@[]^_{|}~"

	var encoded string
	for timestamp > 0 {
		remainder := timestamp % base
		encoded = string(chars[remainder]) + encoded
		timestamp = timestamp / base
	}
	aux := encoded
	if len(encoded) > 5 {
		encoded = encoded[len(encoded)-5:] // Tomar los últimos 5 caracteres
	} else if len(encoded) < 5 {
		// Rellenar con caracteres a la izquierda si la longitud es menor
		for len(encoded) < 5 {
			encoded = "0" + encoded
		}
	}
	helperLog.Logger.Info().Msgf("URL: %s -> %s", aux, encoded)
	return encoded
}
