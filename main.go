package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shortener-url/backend/controller"
	"github.com/shortener-url/backend/repository"
	"github.com/shortener-url/backend/service"
	"gorm.io/gorm"

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
	// Crear el esquema si no existe
	if err := createSchema(db, os.Getenv("DB_POSTGRES_SCHEMA")); err != nil {
		helperLog.Logger.Panic().AnErr("error al crear el esquema!", err)
	}
	// AutoMigrate ejecutará las migraciones automáticamente
	if err := db.AutoMigrate(&repository.URL{}, &repository.Click{}); err != nil {
		helperLog.Logger.Panic().AnErr("error al ejecutar automigrate!", err)
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
	group.GET("/:key-short", controller.GetUrl)
	group.GET("/:key-short/exists", controller.ExistsShortURL)
	group.PATCH("/shorten", controller.UpdateShortURL)
	group.GET("/:key-short/stats", controller.Stats)
	group.GET("/stats", controller.Stats)
	localIP, err := helperCommon.GetLocalIP()
	if err != nil {
		localIP = "localhost"
	}
	swaggerHost := localIP + ":" + port
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("http://"+swaggerHost+"/swagger/doc.json")))
	r.Run()
}

func createSchema(db *gorm.DB, schemaName string) error {
	return db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName).Error
}
