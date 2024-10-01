package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Модель со всей информацией о песне, включая смещение и лимит количества результатов, возвращаемых в этом хэндлере
type queryStringAllSongs struct {
	Group       string `form:"group"`
	Song        string `form:"song"`
	ReleaseDate string `form:"releaseDate,omitempty"`
	Text        string `form:"text,omitempty"`
	Link        string `form:"link,omitempty"`
	Offset      string `form:"offset,omitempty"`
	Limit       string `form:"limit,omitempty"`
}

// Хэндлер для получения песен
func (a *API) GetSongs(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'Get: GetSongs api/songs'")

	// Парсим query string
	var aSongs queryStringAllSongs
	err := c.Bind(&aSongs)
	// Проводим проверки что query string предоставленный пользователем удовлетворяет условиям для данного хэндлера
	if err != nil || aSongs.Offset == "" || aSongs.Limit == "" {
		a.logger.Error("User provide uncorrected query string in url")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: offset and limit value must be not empty"})
		return
	}

	// Считываем значения смещения
	offsetVal, err := strconv.Atoi(aSongs.Offset)
	if err != nil {
		a.logger.Error("User provide uncorrected offset value in url")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: offset value must be a number"})
		return
	}

	// Считываем значения лимита
	limitVal, err := strconv.Atoi(aSongs.Limit)
	if err != nil {
		a.logger.Error("User provide uncorrected limit value in url")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: limit value must be a number"})
		return
	}

	// Далее формируем запрос в БД
	query, err := createQueryDB(aSongs, offsetVal, limitVal)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to generate a query for the DB: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Выполняем запрос в БД
	songs, err := a.storage.Song().GetSongs(query)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table music): %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Если все прошло хорошо, формируем ответ пользователю
	c.JSON(http.StatusOK, ResponceAllSongs{Songs: songs})

	// Логируем окончание запроса
	a.logger.Info("Request 'Get: GetSongs api/songs' successfully done")
}

// Функция для формирования запроса в БД на получение данных библиотеки в зависимости от входящих параметров фильтрации и пагинации
func createQueryDB(song queryStringAllSongs, offset, limit int) (string, error) {
	var queryDB strings.Builder
	var err error
	_, err = queryDB.WriteString(fmt.Sprintf("SELECT * FROM %s ", os.Getenv("TABLE_NAME")))
	if err != nil {
		return "", err
	}
	if song.Group != "" {
		_, err = queryDB.WriteString(fmt.Sprintf("WHERE group=%s ", song.Group))
		if err != nil {
			return "", err
		}
	}
	if song.Song != "" && song.Group != "" {
		_, err = queryDB.WriteString(fmt.Sprintf("& song=%s ", song.Song))
		if err != nil {
			return "", err
		}
	} else if song.Song != "" {
		_, err = queryDB.WriteString(fmt.Sprintf("WHERE song=%s ", song.Song))
		if err != nil {
			return "", err
		}
	}
	if song.ReleaseDate != "" && (song.Group != "" || song.Song != "") {
		_, err = queryDB.WriteString(fmt.Sprintf("& releaseDate=%s ", song.ReleaseDate))
		if err != nil {
			return "", err
		}
	} else if song.ReleaseDate != "" {
		_, err = queryDB.WriteString(fmt.Sprintf("WHERE releaseDate=%s ", song.ReleaseDate))
		if err != nil {
			return "", err
		}
	}
	if song.Text != "" && (song.Group != "" || song.Song != "" || song.ReleaseDate != "") {
		_, err = queryDB.WriteString(fmt.Sprintf("& ARRAY_TO_STRING(text) LIKE '%%%s%%' ", song.Text))
		if err != nil {
			return "", err
		}
	} else if song.Text != "" {
		_, err = queryDB.WriteString(fmt.Sprintf("WHERE ARRAY_TO_STRING(text) LIKE '%%%s%%' ", song.Text))
		if err != nil {
			return "", err
		}
	}
	if song.Link != "" && (song.Group != "" || song.Song != "" || song.ReleaseDate != "" || song.Text != "") {
		_, err = queryDB.WriteString(fmt.Sprintf("& link=%s ", song.Link))
		if err != nil {
			return "", err
		}
	} else if song.Link != "" {
		_, err = queryDB.WriteString(fmt.Sprintf("WHERE link=%s ", song.Link))
		if err != nil {
			return "", err
		}
	}

	_, err = queryDB.WriteString(fmt.Sprintf("OFFSET %v LIMIT %v", offset, limit))
	if err != nil {
		return "", err
	}

	return queryDB.String(), nil
}
