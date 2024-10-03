package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Модель с основной информацией о песне (для работы с query string)
type queryStringSong struct {
	Group string `form:"group"`
	Song  string `form:"song"`
}

// DeleteSong godoc
// @Summary Delete song on given info
// @Produce json
// @Param group path string true "Name of group"
// @Param song path string true "Name of song"
// @Success 200 {object} ResponceMessage
// @Failure 400 {object} ResponceMessage
// @Failure 500 {object} ResponceMessage
// @Router /api/song [delete]

// Хэндлер для удаления песни
func (a *API) DeleteSong(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'DELETE: DeleteSong api/song'")

	// Парсим query string
	var qSong queryStringSong
	err := c.ShouldBindQuery(&qSong)
	// Проверка query string (удовлетворяет ли она условиям данного хэндлера)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with bind query string: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}
	if qSong.Group == "" || qSong.Song == "" {
		a.logger.Error("User provide uncorrected query string in url: group or song is empty")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: group and song value must be not empty"})
		return
	}

	// Логируем обращение к БД
	a.logger.Debug("Sending a request to DB: CheckSong")

	// Ищем песню в БД
	err = a.storage.Song().CheckSong(qSong.Group, qSong.Song)
	// Если песня не найдена
	if err != nil && err == sql.ErrNoRows {
		a.logger.Info(fmt.Sprintf("User trying to delete non existed song. Group: %s, song: %s", qSong.Group, qSong.Song))
		c.JSON(http.StatusBadRequest, ResponceMessage{"You trying to delete non existed song"})
		return
	}
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Логируем обращение к БД
	a.logger.Debug("Sending a request to DB: DeleteSong")

	// Если песня найдена, то удаляем ее
	err = a.storage.Song().DeleteSong(qSong.Group, qSong.Song)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Возвращаем пользователю сообщение об успешно выполненной операции
	c.JSON(http.StatusOK, ResponceMessage{fmt.Sprintf("Song successfully delete. Group: %s, song: %s", qSong.Group, qSong.Song)})

	// Логируем окончание запроса
	a.logger.Info("Request 'DELETE: DeleteSong api/song' successfully done")
}
