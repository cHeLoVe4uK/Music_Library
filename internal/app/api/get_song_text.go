package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Модель с основной информацией о песне и данными о смещении и лимите текста песни, возвращаемого в этом хэндлере
type queryStringSongText struct {
	Group  string `form:"group"`
	Song   string `form:"song"`
	Offset string `form:"offset"`
	Limit  string `form:"limit"`
}

// Хэндлер для получения текста песни
func (a *API) GetSongText(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'GET: GetSongText api/song/text'")

	// Парсим тело запроса
	var song queryStringSongText
	err := c.Bind(&song)
	// Проводим проверки что query string предоставленный пользователем удовлетворяет условиям для данного хэндлера
	if err != nil || song.Group == "" || song.Song == "" || song.Offset == "" || song.Limit == "" {
		a.logger.Error("User provide uncorrected query string in url")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: group, song, offset and limit value must be not empty"})
		return
	}

	// Считываем значения смещения
	offsetVal, err := strconv.Atoi(song.Offset)
	if err != nil {
		a.logger.Error("User provide uncorrected offset value in url")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: offset value must be a number"})
		return
	}

	// Считываем значения лимита
	limitVal, err := strconv.Atoi(song.Limit)
	if err != nil {
		a.logger.Error("User provide uncorrected limit value in url")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: limit value must be a number"})
		return
	}

	// Ищем песню в БД
	ok, err := a.storage.Song().CheckSong(song.Group, song.Song)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table music): %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Если песня не найдена
	if !ok {
		a.logger.Info("User trying to recieve text of non existed song")
		c.JSON(http.StatusBadRequest, ResponceMessage{"You trying to recieve text of non existed song"})
		return
	}

	// Если песня найдена извлекаем текст из БД
	text, err := a.storage.Song().GetSongText(song.Group, song.Song, offsetVal, limitVal)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table music): %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}
	akk := string(text)
	fmt.Println(akk)

	// // Если все прошло хорошо, формируем ответ пользователю
	// c.JSON(http.StatusOK, ResponceTextSong{Text: text})

	// Логируем окончание запроса
	a.logger.Info("Request 'GET: GetSongText api/song/text' successfully done")
}
