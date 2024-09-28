package storage

// Сущность модельного репозитория
type SongRepository struct {
	storage *Storage // Хранит в себе БД, т.к. общение с БД реализовано посредством репозитория (небольшое замыкание)
}

// Метод для получения всех песен
func (song *SongRepository) GetSongs() {}

// Метод для получения текста песни
func (song *SongRepository) GetSongText() {}

// Метод для изменения песни
func (song *SongRepository) UpdateSong() {}

// Метод для удаления песни
func (song *SongRepository) DeleteSong() {}

// Метод для добавления песни
func (song *SongRepository) AddSong() {}

// Метод для получения песни
func (song *SongRepository) GetSong() {}
