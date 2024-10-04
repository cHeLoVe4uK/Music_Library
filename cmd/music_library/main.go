package main

import (
	"log"
	"mus_lib/internal/app/api"
	"mus_lib/migrations"

	"github.com/joho/godotenv"
)

//	@title			Music Library
//	@version		1.0
//	@description	This is a sample server provided CRUD features

//	@host		localhost:8080
//	@BasePath	/api

func init() {
	// Считываем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can't find .env file with configuration")
	}

	// Регистрируем миграции доступные в приложении
	migrations.AddMigration()
}

func main() {
	// Создаем инстанс нашего приложения (сервера)
	server := api.New()

	// Конфигурируем его
	err := server.ConfigureServer()
	if err != nil {
		log.Fatalf("An error occured while configure server: %s", err)
	}

	// Запускаем
	server.StartServer()
}
