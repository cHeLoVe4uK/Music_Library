package api

import "mus_lib/internal/app/models"

// Модель ответа пользователю
type ResponceMessage struct {
	Message string `json:"message"`
}

// Модель ответа пользователю для возвращения текста песни по куплетам
type ResponceTextSong struct {
	Text []string `json:"text"`
}

// Модель ответа пользователю для возвращения песен
type ResponceAllSongs struct {
	Songs []*models.Song `json:"songs"`
}
