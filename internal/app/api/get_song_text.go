package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Модель ответа пользователю для возвращения текста песни по куплетам
type ResponceTextSong struct {
	Verses []string `json:"verses"`
}

// Модель с основной информацией о песне и данными о смещении и лимите куплетов текста песни (для работы с query string)
type queryStringSongText struct {
	Group  string `form:"group"`
	Song   string `form:"song"`
	Offset string `form:"offset"`
	Limit  string `form:"limit"`
}

// GetSongText godoc
// @Summary Retrieve song's text in verses on given info
// @Produce json
// @Param group string true "Name of group"
// @Param song string true "Name of song"
// @Success 200 {object} ResponceTextSong
// @Failure 400 {object} ResponceMessage
// @Failure 500 {object} ResponceMessage
// @Router /api/song/text [get]

// Хэндлер для получения текста песни
func (a *API) GetSongText(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'GET: GetSongText api/song/text'")

	// Парсим query string
	var song queryStringSongText
	err := c.ShouldBindQuery(&song)
	// Проверка query string (удовлетворяет ли она условиям данного хэндлера)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with bind query string: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}
	if song.Group == "" || song.Song == "" || song.Offset == "" || song.Limit == "" {
		a.logger.Error("User provide uncorrected query string in url: group, song, offset or limit is empty")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: group, song, offset and limit value must be not empty"})
		return
	}

	// Считываем значения смещения
	offsetVal, err := strconv.Atoi(song.Offset)
	if err != nil {
		a.logger.Error(fmt.Sprintf("User provide uncorrected offset value in url: %s", err))
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: offset value must be a number"})
		return
	}

	// Считываем значения лимита
	limitVal, err := strconv.Atoi(song.Limit)
	if err != nil {
		a.logger.Error(fmt.Sprintf("User provide uncorrected limit value in url: %s", err))
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: limit value must be a number"})
		return
	}

	// Логируем обращение к БД
	a.logger.Debug("Sending a request to DB: CheckSong")

	// Ищем песню в БД
	err = a.storage.Song().CheckSong(song.Group, song.Song)
	// Если песня не найдена
	if err != nil && err == sql.ErrNoRows {
		a.logger.Info(fmt.Sprintf("User trying to get text of non existed song. Group: %s, song: %s", song.Group, song.Song))
		c.JSON(http.StatusBadRequest, ResponceMessage{"You trying to get text of non existed song"})
		return
	}
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Логируем обращение к БД
	a.logger.Debug("Sending a request to DB: GetSongText")

	// Если песня найдена извлекаем ее текст из БД
	verses, err := a.storage.Song().GetSongText(song.Group, song.Song, offsetVal, limitVal)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Возвращаем пользователю сообщение об успешно выполненной операции (Текст песни по сути будет преобразован в читаемый вид на стороне фронта)
	c.JSON(http.StatusOK, ResponceTextSong{Verses: verses})

	// Логируем окончание запроса
	a.logger.Info("Request 'GET: GetSongText api/song/text' successfully done")
}
