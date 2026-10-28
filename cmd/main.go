package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shortener-url/backend/internal/adapters/handlers"
	"github.com/shortener-url/backend/internal/adapters/repositories"
	"github.com/shortener-url/backend/internal/core/services"
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
	// Note: Using the migration structs from the repository implementation
	if err := db.AutoMigrate(&repositories.PostgresURL{}, &repositories.PostgresClick{}); err != nil {
		helperLog.Logger.Panic().AnErr("error al ejecutar automigrate!", err)
	}
	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err == nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
			defer cancel()
			err = sqlDB.PingContext(ctx)
		}
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "database": "disconnected"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "database": "connected"})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Usa el puerto 8080 por defecto si no está definido
	}

	repo := repositories.NewPostgresRepository(db)
	svc := services.NewURLService(repo)
	handler := handlers.NewHTTPHandler(svc)

	group := r.Group("/api/v1/")
	group.POST("/shorten", handler.Shorten)
	group.GET("/:key-short", handler.GetUrl)
	group.GET("/:key-short/exists", handler.ExistsShortURL)
	group.PATCH("/shorten", handler.UpdateShortURL)
	group.GET("/:key-short/stats", handler.Stats)
	group.GET("/stats", handler.Stats)

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
