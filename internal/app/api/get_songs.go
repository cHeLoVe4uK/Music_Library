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
type responceAllSongs struct {
	Songs []*models.Song `json:"songs"`
}

// Модель со всей информацией о песне, включая смещение и лимит количества возвращаемых результатов из БД (для работы с query string)
type queryStringAllSongs struct {
	Group       string `form:"group"`
	Song        string `form:"song"`
	ReleaseDate string `form:"releaseDate"`
	Text        string `form:"text"`
	Link        string `form:"link"`
	Offset      string `form:"offset"`
	Limit       string `form:"limit"`
}

// GetSongs godoc
//	@Summary		GetSongs
//	@Tags			song
//	@Description	Retrieve songs on given info
//	@Produce		json
//	@Param			offset		path		integer	true	"Offset from the beginning of the list extracted songs"
//	@Param			limit		path		integer	true	"Limit of quantity extracted songs"
//	@Param			group		path		string	"Name of group"
//	@Param			song		path		string	"Name of song"
//	@Param			releaseDate	path		string	"Release date of song"
//	@Param			text		path		string	"Words that will be used to search for songs"
//	@Param			link		path		string	"Link of song on youtube"
//	@Success		200			{array}		models.Song
//	@Failure		400			{object}	responceMessage
//	@Failure		500			{object}	responceMessage
//	@Router			/songs [get]

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
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if aSongs.Offset == "" || aSongs.Limit == "" {
		a.logger.Error("User provide uncorrected query string in url: offset or limit is empty")
		c.JSON(http.StatusBadRequest, errorMessage{"URL have uncorrected parameters in the query string: offset and limit value must be not empty"})
		return
	}

	// Считываем значения смещения
	offsetVal, err := strconv.Atoi(aSongs.Offset)
	if err != nil {
		a.logger.Error(fmt.Sprintf("User provide uncorrected offset value in url: %s", err))
		c.JSON(http.StatusBadRequest, errorMessage{"URL have uncorrected parameters in the query string: offset value must be a number"})
		return
	}

	// Считываем значения лимита
	limitVal, err := strconv.Atoi(aSongs.Limit)
	if err != nil {
		a.logger.Error(fmt.Sprintf("User provide uncorrected limit value in url: %s", err))
		c.JSON(http.StatusBadRequest, errorMessage{"URL have uncorrected parameters in the query string: limit value must be a number"})
		return
	}

	// Формируем запрос в БД
	query, err := createQueryDB(aSongs, offsetVal, limitVal)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to generate a query for the DB: %s", err))
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	// Логируем получившийся запрос (может пригодиться) и обращение к БД
	a.logger.Debug(query)
	a.logger.Debug("Sending a request to DB: GetSongs")

	// Выполняем запрос в БД
	songs, err := a.storage.Song().GetSongs(query)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with connecting to DB (table %s): %s", os.Getenv("TABLE_NAME"), err))
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if len(songs) == 0 {
		a.logger.Info(fmt.Sprintf("No found songs in DB (table %s)", os.Getenv("TABLE_NAME")))
		c.JSON(http.StatusNotFound, errorNotFoundMessage{"No found songs"})
		return
	}

	// Возвращаем пользователю сообщение об успешно выполненной операции (предполагается, что текст песен будет преобразован в читаемый вид на стороне фронта)
	c.JSON(http.StatusOK, responceAllSongs{Songs: songs})

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
