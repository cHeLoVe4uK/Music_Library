package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *API) MockInfo(c *gin.Context) {
	// Парсим query string
	var qSong queryStringSong
	err := c.ShouldBindQuery(&qSong)
	// Проводим проверки что query string предоставленный пользователем удовлетворяет условиям для данного хэндлера
	if err != nil {
		a.logger.Error(fmt.Sprintf("Trouble with bind query string: %s", err))
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if qSong.Group == "" || qSong.Song == "" {
		a.logger.Error("User provide uncorrected query string in url: group or song is empty")
		c.JSON(http.StatusBadRequest, errorMessage{"URL have uncorrected parameters in the query string: group and song value must be not empty"})
		return
	}

	c.JSON(http.StatusOK, externalSong{ReleaseDate: "01.01.1990", Text: "AAAA\nBBBB\nCCCC\n\nDDDD\nEEEE\nFFFF\n\nGGG\nHHH\nKKK", Link: "habdulala"})
}
