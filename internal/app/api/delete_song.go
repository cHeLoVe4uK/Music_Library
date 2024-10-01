package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Модель с основной информацией о песне (для работы с query string)
type queryStringSong struct {
	Group string `form:"group"`
	Song  string `form:"song"`
}

// Хэндлер для удаления песни
func (a *API) DeleteSong(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'DELETE: DeleteSong api/song'")

	// Парсим query string
	var qSong queryStringSong
	err := c.Bind(&qSong)
	// Проводим проверки что query string предоставленный пользователем удовлетворяет условиям для данного хэндлера
	if err != nil || qSong.Group == "" || qSong.Song == "" {
		a.logger.Error("User provide uncorrected query string in url")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: group and song value must be not empty"})
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
		a.logger.Info("User trying to delete non existed song")
		c.JSON(http.StatusBadRequest, ResponceMessage{"You trying to delete non existed song"})
		return
	}

	// Если песня была найдена удаляем ее
	err = a.storage.Song().DeleteSong(qSong.Group, qSong.Song)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table music): %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Если все прошло хорошо, формируем ответ пользователю
	c.JSON(http.StatusOK, ResponceMessage{fmt.Sprintf("Song with name: %s, group: %s successfully delete", qSong.Song, qSong.Group)})

	// Логируем окончание запроса
	a.logger.Info("Request 'DELETE: DeleteSong api/song' successfully done")
}
