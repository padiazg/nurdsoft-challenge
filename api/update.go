package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/padiazg/nurdsoft-challenge/internals"
	"github.com/padiazg/nurdsoft-challenge/models"
)

func updateHandlerFn(s *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			ids        = ctx.Param("id")
			data       = &models.Book{}
			id         int64
			err        error
			book       *models.Book
			statusCode int
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

		book, err = s.data.Update(int32(id), data)
		if err != nil {
			if errors.As(err, &internals.ErrorNotFound{}) {
				statusCode = http.StatusNotFound
			} else {
				statusCode = http.StatusBadRequest
			}

			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, book)
	}
}

func (s *Server) RegisterUpdate() {
	s.router.PUT(s.config.Root+"/:id", updateHandlerFn(s))
}
