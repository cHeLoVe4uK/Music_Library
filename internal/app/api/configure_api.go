package api

import (
	"log/slog"
	"mus_lib/storage"
	"os"

	"github.com/gin-gonic/gin"
)

// Конфигурируем логгер сервера
func (api *API) configureLoggerField() {
	api.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// Конфигурируем роутер сервера
func (api *API) configureRouterField() {
	router := gin.Default()
	apiGroup := router.Group("/api")
	apiGroup.GET("/songs", api.GetSongs)
	apiGroup.GET("/songs/text", api.GetSongText)
	apiGroup.DELETE("/songs", api.DeleteSong)
	apiGroup.PUT("/songs", api.UpdateSong)
	apiGroup.POST("/songs", api.AddSong)

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
