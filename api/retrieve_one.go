package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/padiazg/nurdsoft-challenge/models"
)

func retrieveOneHandlerFn(s *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			ids  = ctx.Param("id")
			err  error
			book *models.Book
			id   int64
		)

		id, err = strconv.ParseInt(ids, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"parsing id to int:": err.Error()})
			return
		}

		book, err = s.data.GetOne(int32(id))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"retrieving data:": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, book)
	}
}

func (s *Server) RegisterRetrieveOne() {
	s.router.GET(s.config.Root+"/:id", retrieveOneHandlerFn(s))
}
