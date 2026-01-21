package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shortener-url/backend/internal/core/domain"
	"github.com/shortener-url/backend/internal/core/ports"
	"github.com/shortener-url/backend/utils"
	"gitlab.com/uchile1/helper/helperHttp"
	"gitlab.com/uchile1/helper/helperLog"
)

type HTTPHandler struct {
	service ports.URLService
}

func NewHTTPHandler(service ports.URLService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

// Shorten godoc
// @Summary      Shorten a URL
// @Description  Creates a shortened URL from a long URL
// @Tags         URL
// @Accept       json
// @Produce      json
// @Param        request body domain.ShortenRequest true "Shorten Request"
// @Success      200 {object} helperHttp.Response
// @Failure      400 {object} helperHttp.Response
// @Failure      500 {object} helperHttp.Response
// @Router       /shorten [post]
func (h *HTTPHandler) Shorten(c *gin.Context) {
	var req domain.ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		helperLog.Logger.Error().Err(err).Msg("invalid request body")
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	_, err := url.ParseRequestURI(req.OriginalURL)
	if err != nil {
		helperLog.Logger.Error().Err(err).Msgf("invalid URL: %s", req.OriginalURL)
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	l, _ := strconv.ParseInt(os.Getenv("KEY_LEN"), 10, 0)
	try, _ := strconv.ParseInt(os.Getenv("KEY_TRY"), 10, 0)
	short, err := h.service.Create(req.OriginalURL, int(try), int(l))
	if err != nil {
		helperLog.Logger.Error().Err(err).Msg("failed to create short URL")
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(short))
}

// GetUrl godoc
// @Summary      Get original URL
// @Description  Redirects or returns the original URL for a given short key
// @Tags         URL
// @Param        key-short path string true "Short Key"
// @Success      200 {object} helperHttp.Response
// @Failure      400 {object} helperHttp.Response
// @Failure      404 {object} helperHttp.Response
// @Failure      500 {object} helperHttp.Response
// @Router       /{key-short} [get]
func (h *HTTPHandler) GetUrl(c *gin.Context) {
	key := c.Param("key-short")
	if key == "" {
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("missing key-short parameter").Error()))
		return
	}
	urlLong, err := h.service.Get(key, c.ClientIP(), c.Request.UserAgent(), c.GetHeader("Referer"))
	if err != nil {
		helperLog.Logger.Error().Err(err).Msg("failed to get URL")
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	if urlLong == nil {
		helperLog.Logger.Error().Err(err).Msg("URL not found")
		c.JSON(http.StatusNotFound, helperHttp.NewErrorResponse("URL not found"))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(urlLong))
}

// ExistsShortURL godoc
// @Summary      Check if short URL exists
// @Description  Verifies if a specific short key is already in use
// @Tags         URL
// @Param        key-short path string true "Short Key"
// @Success      200 {object} helperHttp.Response
// @Failure      400 {object} helperHttp.Response
// @Failure      500 {object} helperHttp.Response
// @Router       /{key-short}/exists [get]
func (h *HTTPHandler) ExistsShortURL(c *gin.Context) {
	key := c.Param("key-short")
	if key == "" {
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("missing key-short parameter").Error()))
		return
	}
	exists, err := h.service.ExistsShortURL(key)
	if err != nil {
		helperLog.Logger.Error().Err(err).Msg("failed to check URL existence")
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(exists))
}

// UpdateShortURL godoc
// @Summary      Update short URL key
// @Description  Changes the short key of an existing URL
// @Tags         URL
// @Accept       json
// @Produce      json
// @Param        request body domain.URLUpdateRequest true "Update Request"
// @Success      200 {object} helperHttp.Response
// @Failure      400 {object} helperHttp.Response
// @Failure      500 {object} helperHttp.Response
// @Router       /shorten [patch]
func (h *HTTPHandler) UpdateShortURL(c *gin.Context) {
	var update domain.URLUpdateRequest
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	if len(update.ShortURLNew) < 2 || len(update.ShortURLNew) > 10 {
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("short URL must be between 2 and 10 characters long").Error()))
		return
	}
	if !utils.IsValidKey(update.ShortURLNew) {
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse("short URL can only contain alphanumeric characters"))
		return
	}
	urlNew, err := h.service.UpdateShortURL(update.ShortURLOld, update.ShortURLNew)
	if err != nil {
		helperLog.Logger.Error().Err(err).Msg("failed to update short URL")
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(urlNew))
}

// Stats godoc
// @Summary      Get URL stats
// @Description  Retrieves statistics (clicks) for a specific short key within a date range
// @Tags         Stats
// @Param        key-short path string false "Short Key (optional for general stats)"
// @Param        date-start query string true "Start Date (YYYY/MM/DD)"
// @Param        date-end query string true "End Date (YYYY/MM/DD)"
// @Success      200 {object} helperHttp.Response
// @Failure      400 {object} helperHttp.Response
// @Failure      500 {object} helperHttp.Response
// @Router       /stats [get]
// @Router       /{key-short}/stats [get]
func (h *HTTPHandler) Stats(c *gin.Context) {
	key := c.Param("key-short")
	startDate, existsStart := c.GetQuery("date-start")
	endDate, existsEnd := c.GetQuery("date-end")
	if !existsStart || !existsEnd {
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("missing date-start or date-end query parameters").Error()))
		return
	}
	layout := "2006/01/02 15:04:05"
	startDateParsed, errStart := time.ParseInLocation(layout, fmt.Sprintf("%s 00:00:00", startDate), time.Local)
	endDateParsed, errEnd := time.ParseInLocation(layout, fmt.Sprintf("%s 23:59:59", endDate), time.Local)
	if errStart != nil || errEnd != nil {
		helperLog.Logger.Error().Err(errStart).Err(errEnd).Msg("invalid date format")
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("invalid date format").Error()))
		return
	}
	const hoursInYear = 366 * 24
	const maxDuration = time.Hour * hoursInYear
	if endDateParsed.Before(startDateParsed) || endDateParsed.After(time.Now().Add(time.Hour*24)) || endDateParsed.Sub(startDateParsed) > maxDuration {
		helperLog.Logger.Error().Msgf("invalid date range start: %s, end: %s", startDate, endDate)
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("invalid date range start: %s, end: %s", startDate, endDate).Error()))
		return
	}
	urls, err := h.service.FindByParams(key, startDateParsed, endDateParsed)
	if err != nil {
		helperLog.Logger.Error().Err(err).Msgf("failed to get URL stats by short URL: %s", key)
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(urls))
}
