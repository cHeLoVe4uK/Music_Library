package api

import (
	"encoding/json"
	"fmt"
	"mus_lib/internal/app/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Модель с основной информацией о песне (для работы с request body)
type requestBodySong struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// Модель с детальной информацией о песне, принимаемой со стороннего API
type externalSong struct {
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

// Хэндлер для добавления песни
func (a *API) AddSong(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'POST: AddSong api/song'")

	// Парсим тело запроса
	var reqSong requestBodySong
	err := c.ShouldBindJSON(&reqSong)
	// Проводим проверки что json предоставленный пользователем удовлетворяет условиям для данного хэндлера
	if err != nil || reqSong.Group == "" || reqSong.Song == "" {
		a.logger.Error("User provide uncorrected JSON")
		c.JSON(http.StatusBadRequest, ResponceMessage{Message: "You provide uncorrected JSON"})
		return
	}

	// Обращаемся к стороннему сервису для получения данных о песне (предполагается, что на нем реализована вся логика ко входящему запросу)
	resp, err := http.Get(fmt.Sprintf("https://api.example.com/info?group=%s&song=%s", reqSong.Group, reqSong.Song))
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to fetch song data: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{Message: "Server error. Try later"})
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа со стороннего сервера (если он равен 400, то скорее всего пользователь предоставил данные несуществующей песни)
	if resp.StatusCode == http.StatusBadRequest {
		a.logger.Error("Failed to fetch song data: song does not exist")
		c.JSON(http.StatusBadRequest, ResponceMessage{"Song does not exist. Check the correctnes of the provided data"})
		return
	}

	// Проверяем статус ответа со стороннего сервера (если он равен 500, то на их стороне какая-то ошибка с сервером)
	if resp.StatusCode == http.StatusInternalServerError {
		a.logger.Error("Failed to fetch song data: server error on the external API side")
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Считываем данные из ответа, полученного со стороннего сервиса и парсим тело ответа в нашу песню
	var extSong externalSong
	err = json.NewDecoder(resp.Body).Decode(&extSong)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to read or unmarshal response body: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Добавляем песню в нашу БД
	song := models.Song{Group: reqSong.Group, Song: reqSong.Song, ReleaseDate: extSong.ReleaseDate, Text: strings.Split(extSong.Text, "\n\n"), Link: extSong.Link}
	err = a.storage.Song().AddSong(&song)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table music): %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Если все прошло хорошо, формируем ответ пользователю
	c.JSON(http.StatusCreated, ResponceMessage{fmt.Sprintf("Song with name: %s, group: %s successfully add", song.Song, song.Group)})

	// Логируем окончание запроса
	a.logger.Info("Request 'POST: AddSong api/song' successfully done")
}
