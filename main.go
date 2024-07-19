package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shortener-url/backend/controller"
	"github.com/shortener-url/backend/repository"
	"github.com/shortener-url/backend/service"

	_ "github.com/shortener-url/backend/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/uchile1/helper/helperCommon"
	"gitlab.com/uchile1/helper/helperDB"
	"gitlab.com/uchile1/helper/helperLog"
	"gitlab.com/uchile1/helper/helperMiddleware"
)

// @title Shortener URL API
// @version 1.0
// @description acortador de URLs
// @host localhost:8080
// @BasePath /api/v1

func main() {
	r := gin.Default()
	r.Use(helperMiddleware.Cors())
	db, err := helperDB.GetDBPostgres()
	if err != nil {
		helperLog.Logger.Panic().AnErr("NO se pudo conectar a la BD!", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Usa el puerto 8080 por defecto si no está definido
	}
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	controller := controller.NewController(service)
	group := r.Group("/api/v1/")
	group.POST("/shorten", controller.Shorten)
	group.GET("/:short-url", controller.GetUrl)
	group.GET("/:short-url/stats", controller.Stats)
	group.GET("/stats", controller.Stats)
	localIP, err := helperCommon.GetLocalIP()
	if err != nil {
		localIP = "localhost"
	}
	swaggerHost := localIP + ":" + port
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("http://"+swaggerHost+"/swagger/doc.json")))
	r.Run()
}
