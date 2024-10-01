package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Хэндлер для изменения песни
func (a *API) UpdateSong(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'PUT: UpdateSong api/song'")

	// Парсим query string
	var qSong queryStringSong
	err := c.Bind(&qSong)
	// Проводим проверки что query string предоставленный пользователем удовлетворяет условиям для данного хэндлера
	if err != nil || qSong.Group == "" || qSong.Song == "" {
		a.logger.Error("User provide uncorrected query string in url")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: group and song value must be not empty"})
		return
	}

	// Парсим тело запроса
	var reqSong requestBodySong
	err = c.ShouldBindJSON(&reqSong)
	// Проводим проверки что json предоставленный пользователем удовлетворяет условиям для данного хэндлера
	if err != nil || reqSong.Group == "" || reqSong.Song == "" {
		a.logger.Error("User provide uncorrected JSON")
		c.JSON(http.StatusBadRequest, ResponceMessage{"You provide uncorrected JSON"})
		return
	}

	// Ищем песню в БД
	ok, err := a.storage.Song().CheckSong(qSong.Group, qSong.Song)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table music): %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Если песня не найдена
	if !ok {
		a.logger.Info("User trying to update non existed song")
		c.JSON(http.StatusBadRequest, ResponceMessage{"You trying to update non existed song"})
		return
	}

	// Обновляем данные песни
	err = a.storage.Song().UpdateSong(reqSong.Group, reqSong.Song, qSong.Group, qSong.Song)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table music): %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Если все прошло хорошо, формируем ответ пользователю
	c.JSON(http.StatusOK, ResponceMessage{fmt.Sprintf("Song with name: %s, group: %s successfully update", qSong.Song, qSong.Group)})

	// Логируем окончание запроса
	a.logger.Info("Request 'PUT: UpdateSong api/song' successfully done")
}
