package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mus_lib/internal/app/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Модель ответа пользователю
type ResponceMessage struct {
	Message string `json:"message"`
}

// Модель с основной информацией о песне (для работы с request body)
type requestBodySong struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// Модель с детальной информацией о песне, принимаемой со стороннего API (для работы с responce body)
type externalSong struct {
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

// Хэндлер для добавления песни
func (a *API) AddSong(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'POST: AddSong api/song'")

	// Парсим request body
	var reqSong requestBodySong
	err := c.ShouldBindJSON(&reqSong)
	// Проверка request body (удовлетворяет ли оно условиям для данного хэндлера)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with bind request body: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}
	if reqSong.Group == "" || reqSong.Song == "" {
		a.logger.Error("User provide uncorrected JSON: group or song is empty")
		c.JSON(http.StatusBadRequest, ResponceMessage{"You provide uncorrected JSON: group and song value must be not empty"})
		return
	}

	// Логируем обращение к БД
	a.logger.Debug("Sending a request to DB: CheckSong")

	// Ищем песню в БД
	err = a.storage.Song().CheckSong(reqSong.Group, reqSong.Song)

	if err != nil && err != sql.ErrNoRows {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}
	// Если песня найдена
	if err == nil {
		a.logger.Info(fmt.Sprintf("User trying to add existed song. Group: %s, song: %s", reqSong.Group, reqSong.Song))
		c.JSON(http.StatusBadRequest, ResponceMessage{"You trying to add existed song"})
		return
	}

	// Логируем обращение к стороннему API
	// a.logger.Debug("Sending a request to external API")
	// Если не найдена обращаемся к стороннему сервису для получения данных о песне (предполагается, что на нем реализована вся логика ко входящему запросу) (это обращение закомментировано, потому что в данном случае это фейк, но просто чтобы было понятно как оно выглядит)
	// resp, err := http.Get(fmt.Sprintf("http://api.example.com/info?group=%s&song=%s", strings.Replace(reqSong.Group, " ", "+", -1), strings.Replace(reqSong.Song, " ", "+", -1)))
	// if err != nil {
	// 	a.logger.Error(fmt.Sprintf("Failed to fetch song data: %s", err))
	// 	c.JSON(http.StatusInternalServerError, ResponceMessage{Message: "Server error. Try later"})
	// 	return
	// }
	// defer resp.Body.Close()

	// Логируем обращение к стороннему API
	a.logger.Debug("Sending a request to external API. Method: Get, path: localhost:8080/info")

	// Это mock обращение, чтбы сымитировать сторонний API
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", strings.Replace(reqSong.Group, " ", "+", -1), strings.Replace(reqSong.Song, " ", "+", -1)))
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to fetch song data: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{Message: "Server error. Try later"})
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа со стороннего сервера (если он равен 400, то скорее всего пользователь предоставил данные несуществующей песни, если он равен 500, то на их стороне какая-то ошибка с сервером) (mock обращение может выдать такие статусы, но эта проверка так же будет действительна и для настоящего стороннего сервера)
	if resp.StatusCode == http.StatusBadRequest {
		a.logger.Error("Failed to fetch song data: song does not exist")
		c.JSON(http.StatusBadRequest, ResponceMessage{"Song does not exist. Check the correctnes of the provided data"})
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		a.logger.Error("Failed to fetch song data: server error on the external API side")
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Парсим responce body, полученный из стороннего сервиса в нашу песню
	var extSong externalSong
	err = json.NewDecoder(resp.Body).Decode(&extSong)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to read response body: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Создаем песню, которую будем добавлять в БД, из полученных данных
	song := models.Song{Group: reqSong.Group, Song: reqSong.Song, ReleaseDate: extSong.ReleaseDate, Text: strings.Split(extSong.Text, "\n\n"), Link: extSong.Link}

	// Логируем обращение к БД
	a.logger.Debug("Sending a request to DB: AddSong")

	// Добавляем песню в БД
	err = a.storage.Song().AddSong(&song)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Возвращаем пользователю сообщение об успешно выполненной операции
	c.JSON(http.StatusCreated, ResponceMessage{fmt.Sprintf("Song successfully add. Group: %s, song: %s", song.Group, song.Song)})

	// Логируем окончание запроса
	a.logger.Info("Request 'POST: AddSong api/song' successfully done")
}
