package storage

import (
	"fmt"
	"mus_lib/internal/app/models"
	"os"
	"strings"

	"github.com/lib/pq"
)

// Сущность модельного репозитория
type SongRepository struct {
	storage *Storage // Хранит в себе БД, т.к. общение с БД реализовано посредством репозитория (небольшое замыкание)
}

// Метод для получения всех песен из БД
func (s *SongRepository) GetSongs(query string) ([]*models.Song, error) {
	res, err := s.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	songs := make([]*models.Song, 0)

	for res.Next() {
		song := models.Song{}
		err := res.Scan(&song.Group, &song.Song, &song.ReleaseDate, pq.Array(&song.Text), &song.Link)
		if err != nil {
			continue
		}

		songs = append(songs, &song)
	}

	return songs, nil
}

// Метод для получения текста песни из БД
func (s *SongRepository) GetSongText(group, song string, offset, limit int) ([]string, error) {
	query := fmt.Sprintf(`SELECT text[$1:$2] FROM %s WHERE "group"=$3 AND song=$4`, os.Getenv("TABLE_NAME"))
	res := s.storage.db.QueryRow(query, offset+1, limit+offset, strings.ToLower(group), strings.ToLower(song))

	var text []string

	err := res.Scan(pq.Array(&text))
	if err != nil {
		return nil, err
	}

	return text, nil
}

// Метод для изменения песни в БД
func (s *SongRepository) UpdateSong(newgroup, newsong, oldgroup, oldsong string) error {
	query := fmt.Sprintf(`UPDATE %s SET "group"=$1, song=$2 WHERE "group"=$3 AND song=$4`, os.Getenv("TABLE_NAME"))

	_, err := s.storage.db.Exec(query, strings.ToLower(newgroup), strings.ToLower(newsong), strings.ToLower(oldgroup), strings.ToLower(oldsong))
	return err
}

// Метод для удаления песни из БД
func (s *SongRepository) DeleteSong(group string, song string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE "group"=$1 AND song=$2`, os.Getenv("TABLE_NAME"))

	_, err := s.storage.db.Exec(query, strings.ToLower(group), strings.ToLower(song))
	return err
}

// Метод для добавления песни в БД
func (s *SongRepository) AddSong(song *models.Song) error {
	query := fmt.Sprintf(`INSERT INTO %s VALUES ($1, $2, $3, $4, $5)`, os.Getenv("TABLE_NAME"))

	_, err := s.storage.db.Exec(query, strings.ToLower(song.Group), strings.ToLower(song.Song), song.ReleaseDate, pq.Array(song.Text), song.Link)
	return err
}

// Метод для проверки наличия песни в БД
func (s *SongRepository) CheckSong(group string, song string) error {
	query := fmt.Sprintf(`SELECT "group", song FROM %s WHERE "group"=$1 AND song=$2`, os.Getenv("TABLE_NAME"))
	res := s.storage.db.QueryRow(query, strings.ToLower(group), strings.ToLower(song))

	songTemp := models.Song{}

	err := res.Scan(&songTemp.Group, &songTemp.Song)
	return err
}
