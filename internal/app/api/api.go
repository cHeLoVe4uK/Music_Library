package api

import (
	"log"
	"log/slog"
	"mus_lib/storage"

	"os"

	"github.com/gin-gonic/gin"
)

// Инстанс нашего сервера
type API struct {
	// Поля неэкспортируемые (конфендициальная информация)
	logger  *slog.Logger     // логер который будет использоваться в процессе работы сервера
	router  *gin.Engine      // роутер который будет использоваться в процессе работы сервера (в нашем случае используем фреймворк gin)
	storage *storage.Storage // БД, которая будет использоваться в процессе работы сервера
}

// Конструктор, возвращающий инстанс нашего сервера
func New() *API {
	return &API{}
}

// Метод, настраивающий наш сервер
func (api *API) ConfigureServer() error {
	// Настройка поля логгер
	api.configureLoggerField()
	api.logger.Info("Logger succsessfully configured")

	// Настройка поля роутер
	api.configureRouterField()
	api.logger.Info("Router succsessfully configured")

	// Настройка поля с хранилищем
	err := api.configureStorageField()
	if err != nil {
		return err
	}
	api.logger.Info("DB connection succsessfully installed")

	// Сигнал о том, что настройка прошла успешно
	api.logger.Info("Ready to start on port:" + os.Getenv("BIND_ADDR"))

	return nil
}

// Метод, запускающий сервер
func (api *API) StartServer() {
	// Запускаем сервер
	err := api.router.Run(":" + os.Getenv("BIND_ADDR"))
	if err != nil {
		log.Fatal("Server work is over because: ", err)
	}
}
