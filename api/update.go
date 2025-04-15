package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/padiazg/nurdsoft-challenge/models"
)

func updateHandlerFn(s *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			ids  = ctx.Param("id")
			data = &models.Book{}
			id   int64
			err  error
			book *models.Book
		)

		id, err = strconv.ParseInt(ids, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"parsing id to int:": err.Error()})
			return
		}

		if err = ctx.ShouldBindJSON(data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"parsing json": err.Error()})
			return
		}

		if book, err = s.data.Update(int32(id), data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, book)
	}
}

func (s *Server) RegisterUpdate() {
	s.router.PUT(s.config.Root+"/:id", updateHandlerFn(s))
}
