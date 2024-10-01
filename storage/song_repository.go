package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"mus_lib/internal/app/models"
	"os"
	"strings"
)

// Сущность модельного репозитория
type SongRepository struct {
	storage *Storage // Хранит в себе БД, т.к. общение с БД реализовано посредством репозитория (небольшое замыкание)
}

// Метод для получения всех песен
func (song *SongRepository) GetSongs(query string) ([]*models.Song, error) {
	res, err := song.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	songs := make([]*models.Song, 0)

	for res.Next() {
		s := models.Song{}
		err := res.Scan(&s.Group, &s.Song, &s.ReleaseDate, &s.Text, &s.Link)
		if err != nil {
			continue
		}
		songs = append(songs, &s)
	}

	return songs, nil
}

// Метод для получения текста песни
func (song *SongRepository) GetSongText(g, s string, offset, limit int) ([]uint8, error) {
	query := fmt.Sprintf("SELECT text::text[] FROM %s WHERE \"group\"=$1 AND song=$2", os.Getenv("TABLE_NAME"))
	res := song.storage.db.QueryRow(query, strings.ToLower(g), strings.ToLower(s))

	var text []uint8

	err := res.Scan(&text)
	if err != nil {
		fmt.Println("error looks like")
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("here")
			fmt.Println(err)
			return nil, sql.ErrNoRows
		}
		fmt.Println("or here")
		fmt.Println(err)
		return nil, err
	}

	return text, nil
}

// Метод для изменения песни
func (song *SongRepository) UpdateSong(newgroup, newsong, g, s string) error {
	query := fmt.Sprintf("UPDATE %s SET group=$1, song=$2 WHERE group=$3 & song=$4", os.Getenv("TABLE_NAME"))

	_, err := song.storage.db.Exec(query, strings.ToLower(newgroup), strings.ToLower(newsong), strings.ToLower(g), strings.ToLower(s))
	return err
}

// Метод для удаления песни
func (song *SongRepository) DeleteSong(g string, s string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE group=$1 & song=$2", os.Getenv("TABLE_NAME"))

	_, err := song.storage.db.Exec(query, strings.ToLower(g), strings.ToLower(s))
	return err
}

// Метод для добавления песни
func (song *SongRepository) AddSong(s *models.Song) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3, $4, $5)", os.Getenv("TABLE_NAME"))

	_, err := song.storage.db.Exec(query, strings.ToLower(s.Group), strings.ToLower(s.Song), s.ReleaseDate, s.Text, s.Link)
	return err
}

// Метод для проверки наличия песни
func (song *SongRepository) CheckSong(g string, s string) (bool, error) {
	query := fmt.Sprintf("SELECT \"group\", song FROM %s WHERE \"group\"=$1 AND song=$2", os.Getenv("TABLE_NAME"))
	res := song.storage.db.QueryRow(query, strings.ToLower(g), strings.ToLower(s))

	var founded bool
	songTemp := models.Song{}

	err := res.Scan(&songTemp.Group, &songTemp.Song)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return founded, sql.ErrNoRows
		}
		return founded, err
	}

	founded = true
	return founded, nil
}
