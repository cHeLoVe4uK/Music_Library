package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// UpdateSong godoc
// @Summary Update song on given info
// @Accept json
// @Produce json
// @Param group path string true "Name of group"
// @Param song path string true "Name of song"
// @Param group body string true "Name of group"
// @Param song body string true "Name of song"
// @Success 200 {object} ResponceMessage
// @Failure 400 {object} ResponceMessage
// @Failure 500 {object} ResponceMessage
// @Router /api/song [put]

// Хэндлер для изменения песни
func (a *API) UpdateSong(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'PUT: UpdateSong api/song'")

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

	// Парсим request body
	var reqSong requestBodySong
	err = c.ShouldBindJSON(&reqSong)
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
	err = a.storage.Song().CheckSong(qSong.Group, qSong.Song)
	// Если песня не найдена
	if err != nil && err == sql.ErrNoRows {
		a.logger.Info(fmt.Sprintf("User trying to update non existed song. Group: %s, song: %s", qSong.Group, qSong.Song))
		c.JSON(http.StatusBadRequest, ResponceMessage{"You trying to update non existed song"})
		return
	}
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Логируем обращение к БД
	a.logger.Debug("Sending a request to DB: UpdateSong")

	// Если песня найдена, то обновляем ее данные
	err = a.storage.Song().UpdateSong(reqSong.Group, reqSong.Song, qSong.Group, qSong.Song)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Возвращаем пользователю сообщение об успешно выполненной операции
	c.JSON(http.StatusOK, ResponceMessage{fmt.Sprintf("Song successfully update. Group: %s, song: %s", qSong.Group, qSong.Song)})

	// Логируем окончание запроса
	a.logger.Info("Request 'PUT: UpdateSong api/song' successfully done")
}
