package api

import (
	"fmt"
	"mus_lib/internal/app/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Модель ответа пользователю для возвращения песен
type ResponceAllSongs struct {
	Songs []*models.Song `json:"songs"`
}

// Модель со всей информацией о песне, включая смещение и лимит количества возвращаемых результатов из БД (для работы с query string)
type queryStringAllSongs struct {
	Group       string `form:"group"`
	Song        string `form:"song"`
	ReleaseDate string `form:"releaseDate,omitempty"`
	Text        string `form:"text,omitempty"`
	Link        string `form:"link,omitempty"`
	Offset      string `form:"offset,omitempty"`
	Limit       string `form:"limit,omitempty"`
}

// GetSongs godoc
// @Summary Retrieve songs in verses on given info
// @Produce json
// @Param offset path integer true "Offset from the beginning of the list extracted songs"
// @Param limit path integer true "Limit of quantity extracted songs"
// @Param releaseDate path string "Release date of song"
// @Param text path string "Words that will be used to search for songs"
// @Param link path string "Link of song on youtube"
// @Success 200 {object} models.Song
// @Failure 400 {object} ResponceMessage
// @Failure 500 {object} ResponceMessage
// @Router /api/songs [get]

// Хэндлер для получения песен
func (a *API) GetSongs(c *gin.Context) {
	// Логируем начало выполнение запроса
	a.logger.Info("User do 'Get: GetSongs api/songs'")

	// Парсим query string
	var aSongs queryStringAllSongs
	err := c.ShouldBindQuery(&aSongs)
	// Проверка query string (удовлетворяет ли она условиям данного хэндлера)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with bind query string: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}
	if aSongs.Offset == "" || aSongs.Limit == "" {
		a.logger.Error("User provide uncorrected query string in url: offset or limit is empty")
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: offset and limit value must be not empty"})
		return
	}

	// Считываем значения смещения
	offsetVal, err := strconv.Atoi(aSongs.Offset)
	if err != nil {
		a.logger.Error(fmt.Sprintf("User provide uncorrected offset value in url: %s", err))
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: offset value must be a number"})
		return
	}

	// Считываем значения лимита
	limitVal, err := strconv.Atoi(aSongs.Limit)
	if err != nil {
		a.logger.Error(fmt.Sprintf("User provide uncorrected limit value in url: %s", err))
		c.JSON(http.StatusBadRequest, ResponceMessage{"URL have uncorrected parameters in the query string: limit value must be a number"})
		return
	}

	// Формируем запрос в БД
	query, err := createQueryDB(aSongs, offsetVal, limitVal)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to generate a query for the DB: %s", err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Логируем получившийся запрос (может пригодиться) и обращение к БД
	a.logger.Debug(query)
	a.logger.Debug("Sending a request to DB: GetSongs")

	// Выполняем запрос в БД
	songs, err := a.storage.Song().GetSongs(query)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, ResponceMessage{"Server error. Try later"})
		return
	}

	// Возвращаем пользователю сообщение об успешно выполненной операции (Текст песни по сути будет преобразован в читаемый вид на стороне фронта)
	c.JSON(http.StatusOK, ResponceAllSongs{Songs: songs})

	// Логируем окончание запроса
	a.logger.Info("Request 'Get: GetSongs api/songs' successfully done")
}

// Функция для формирования запроса в БД на получение данных библиотеки в зависимости от входящих параметров фильтрации и пагинации
func createQueryDB(song queryStringAllSongs, offset, limit int) (string, error) {
	// Подготавливаем необходимые переменные
	var (
		queryDB          strings.Builder
		err              error
		whereExpressions []string
	)

	// Формируем основу запроса для БД
	_, err = queryDB.WriteString(fmt.Sprintf(`SELECT "group", song, releaseDate, string_to_array(array_to_string(text, E'\n\n'), E'\n\n') as text, link FROM %s `, os.Getenv("TABLE_NAME")))
	if err != nil {
		return "", err
	}

	// Считываем параметры, присутствующие в query string
	if song.Group != "" {
		whereExpressions = append(whereExpressions, `"group"=`+song.Group)
	}
	if song.Song != "" {
		whereExpressions = append(whereExpressions, "song="+song.Song)
	}
	if song.ReleaseDate != "" {
		whereExpressions = append(whereExpressions, "releaseDate="+song.ReleaseDate)
	}
	if song.Text != "" {
		whereExpressions = append(whereExpressions, fmt.Sprintf("array_to_string(text, '') LIKE '%%%s%%' ", song.Text))
	}
	if song.Link != "" {
		whereExpressions = append(whereExpressions, "releaseDate="+song.Link)
	}

	// Проходимся по этим параметрам, чтобы правильно сформировать запрос
	for i := range whereExpressions {
		whereExpression := whereExpressions[i]
		if i == 0 {
			whereExpression = "WHERE " + whereExpression
		} else {
			whereExpression = "AND  " + whereExpression
		}
		_, err = queryDB.WriteString(whereExpression)
		if err != nil {
			return "", err
		}
	}

	// Заканчиваем формирование запроса
	_, err = queryDB.WriteString(fmt.Sprintf("OFFSET %v LIMIT %v", offset, limit))
	if err != nil {
		return "", err
	}

	// Возвращаем запрос
	return queryDB.String(), nil
}
