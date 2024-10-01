package storage

import (
	"context"
	"database/sql"
	"fmt"
	"mus_lib/migrations"
	"os"

	_ "github.com/lib/pq"
)

// Инстанс хранилища для приложения
type Storage struct {
	// Поля неэкспортируемые (конфендициальная информация)
	db             *sql.DB         // Сущность, представляющая собой мост между нашим приложением и БД
	songRepository *SongRepository // Модельный репозиторий, через который будет проводиться работа с БД
}

// Конструктор, возвращающий инстанс нашего хранилища
func New() *Storage {
	return &Storage{}
}

// Метод, открывающий соединение между нашим приложением и БД
func (storage *Storage) Open() error {
	db, err := sql.Open(os.Getenv("DRIVER_NAME"), fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	storage.db = db

	return nil
}

// Метод, закрывающий наше соединение с БД
func (storage *Storage) Close() {
	storage.db.Close()
}

// Метод, создающий таблицу в нашей БД (накатывающий миграцию)
func (storage *Storage) CreateTable() error {
	ctx := context.Background()

	err := migrations.Up(ctx, storage.db)
	return err
}

// Метод, создающий публичный репозиторий для Song
func (storage *Storage) Song() *SongRepository {
	if storage.songRepository != nil {
		return storage.songRepository
	}

	storage.songRepository = &SongRepository{
		storage: storage,
	}

	return storage.songRepository
}
