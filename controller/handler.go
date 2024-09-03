package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shortener-url/backend/model"
	"github.com/shortener-url/backend/service"
	"gitlab.com/uchile1/helper/helperHttp"
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
	short, err := controller.Service.Create(req.OriginalURL, int(try), int(l))
	if err != nil {
		helperLog.Logger.Error().Err(err).Msg("failed to create short URL")
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(short))
}

func (controller *Controller) GetUrl(c *gin.Context) {
	key := c.Param("key-short")
	if key == "" {
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("missing key-short parameter").Error()))
		return
	}
	urlLong, err := controller.Service.Get(key, c.ClientIP(), c.Request.UserAgent(), c.GetHeader("Referer"))
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

func (controller *Controller) ExistsShortURL(c *gin.Context) {
	key := c.Param("key-short")
	if key == "" {
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("missing key-short parameter").Error()))
		return
	}
	exists, err := controller.Service.ExistsShortURL(key)
	if err != nil {
		helperLog.Logger.Error().Err(err).Msg("failed to check URL existence")
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(exists))
}

func (controller *Controller) UpdateShortURL(c *gin.Context) {
	type URLUpdateRequest struct {
		ShortURLNew string `json:"short_url_new" binding:"required"`
		ShortURLOld string `json:"short_url_old" binding:"required"`
	}
	var update URLUpdateRequest
	c.ShouldBindBodyWithJSON(&update)
	urlNew, err := controller.Service.UpdateShortURL(update.ShortURLOld, update.ShortURLNew)
	if err != nil {
		helperLog.Logger.Error().Err(err).Msg("failed to update short URL")
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(urlNew))
}

func (controller *Controller) Stats(c *gin.Context) {
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
	const hoursInYear = 366 * 24 // horas en un año (bisiesto)
	const maxDuration = time.Hour * hoursInYear
	if endDateParsed.Before(startDateParsed) || endDateParsed.After(time.Now().Add(time.Hour*24)) || endDateParsed.Sub(startDateParsed) > maxDuration {
		helperLog.Logger.Error().Msgf("invalid date range start: %s, end: %s", startDate, endDate)
		c.JSON(http.StatusBadRequest, helperHttp.NewErrorResponse(fmt.Errorf("invalid date range start: %s, end: %s", startDate, endDate).Error()))
		return
	}
	urls, err := controller.Service.FindByParams(key, startDateParsed, endDateParsed)
	if err != nil {
		helperLog.Logger.Error().Err(err).Msgf("failed to get URL stats by short URL: %s", key)
		c.JSON(http.StatusInternalServerError, helperHttp.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helperHttp.NewSuccessResponse(urls))
}
