package api

import (
	"log/slog"
	"mus_lib/storage"
	"net/http"
	"os"
	"time"

	_ "mus_lib/docs"

	"github.com/gin-gonic/gin"
	swaggoFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Конфигурируем логгер сервера
func (api *API) configureLoggerField() {
	api.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// Конфигурируем клиент сервера, через который будем обращаться к стороннему API
func (api *API) configureClientField() {
	http.DefaultClient = &http.Client{Timeout: 5 * time.Second}
	api.client = http.DefaultClient
}

// Конфигурируем роутер сервера
func (api *API) configureRouterField() {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggoFiles.Handler))

	apiGroup := router.Group("/api")
	apiGroup.GET("/songs", api.GetSongs)
	apiGroup.GET("/song/text", api.GetSongText)
	apiGroup.DELETE("/song", api.DeleteSong)
	apiGroup.PUT("/song", api.UpdateSong)
	apiGroup.POST("/song", api.AddSong)
	apiGroup.GET("/info", api.MockInfo)

	api.router = router
}

// Конфигурируем хранилище сервера и создаем в нем таблицу
func (api *API) configureStorageField() error {
	storage := storage.New()

	err := storage.Open()
	if err != nil {
		return err
	}

	err = storage.CreateTable()
	if err != nil {
		return err
	}

	api.storage = storage
	return nil
}
