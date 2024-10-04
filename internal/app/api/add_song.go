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

var (
	serverError = errorServerMessage{Message: "Server error. Try later"}
)

// Модель ответа пользователю в случае успешного выполнения хэндлера
type responceMessage struct {
	Message string `json:"message"`
}

// Модель ответа пользователю для указания ошибки
type errorMessage struct {
	Message string `json:"message"`
}

// Модель ответа пользователю для указания ошибки, что песня не найдена
type errorNotFoundMessage struct {
	Message string `json:"message"`
}

// Модель ответа пользователю для указания ошибки сервера
type errorServerMessage struct {
	Message string `json:"message"`
}

// Модель с основной информацией о песне (для работы с request body)
type requestBodySong struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// Модель с детальной информацией о песне, принимаемой со стороннего API (для работы с responce body)
type externalSong struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// AddSong godoc
//	@Summary		AddSong
//	@Tags			song
//	@Description	Create song on given info
//	@Accept			json
//	@Produce		json
//	@Param			input	body		requestBodySong	true	"Song info"
//	@Success		201		{object}	responceMessage
//	@Failure		400		{object}	responceMessage
//	@Failure		500		{object}	responceMessage
//	@Router			/song [post]

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
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if reqSong.Group == "" || reqSong.Song == "" {
		a.logger.Error("User provide uncorrected JSON: group or song is empty")
		c.JSON(http.StatusBadRequest, errorMessage{"You provide uncorrected JSON: group and song value must be not empty"})
		return
	}

	// Логируем обращение к БД
	a.logger.Debug("Sending a request to DB: CheckSong")

	// Ищем песню в БД
	err = a.storage.Song().CheckSong(reqSong.Group, reqSong.Song)

	if err != nil && err != sql.ErrNoRows {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	// Если песня найдена
	if err == nil {
		a.logger.Info(fmt.Sprintf("User trying to add existed song. Group: %s, song: %s", reqSong.Group, reqSong.Song))
		c.JSON(http.StatusBadRequest, errorMessage{"You trying to add existed song"})
		return
	}

	// Логируем обращение к стороннему API (mock обращение)
	a.logger.Debug("Sending a request to external API. Method: Get, path: localhost:8080/info")

	// Это mock обращение, чтобы сымитировать сторонний API (осуществляется 3 раза в случае истечения таймаута)
	var responce *http.Response
	for i := 0; i < 3; i++ {
		resp, err := a.client.Get(fmt.Sprintf("http://localhost:8080/api/info?group=%s&song=%s", strings.Replace(reqSong.Group, " ", "+", -1), strings.Replace(reqSong.Song, " ", "+", -1)))
		if err == nil {
			responce = resp
			break
		}
		if i == 2 {
			a.logger.Error(fmt.Sprintf("Failed to fetch song data: %s", err))
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	defer responce.Body.Close()

	// Проверяем статус ответа со стороннего сервера (если он равен 400, то скорее всего пользователь предоставил данные несуществующей песни, если он равен 500, то на их стороне какая-то ошибка с сервером) (mock обращение может выдать такие статусы, но эта проверка так же будет действительна и для настоящего стороннего сервера)
	if responce.StatusCode == http.StatusBadRequest {
		a.logger.Error("Failed to fetch song data: song does not exist")
		c.JSON(http.StatusBadRequest, errorMessage{"Song does not exist. Check the correctnes of the provided data"})
		return
	}
	if responce.StatusCode == http.StatusInternalServerError {
		a.logger.Error("Failed to fetch song data: server error on the external API side")
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	// Парсим responce body, полученный из стороннего сервиса в нашу песню
	var extSong externalSong
	err = json.NewDecoder(responce.Body).Decode(&extSong)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to read response body: %s", err))
		c.JSON(http.StatusInternalServerError, serverError)
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
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	// Возвращаем пользователю сообщение об успешно выполненной операции
	c.JSON(http.StatusCreated, responceMessage{fmt.Sprintf("Song successfully add. Group: %s, song: %s", song.Group, song.Song)})

	// Логируем окончание запроса
	a.logger.Info("Request 'POST: AddSong api/song' successfully done")
}
