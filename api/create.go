package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/padiazg/nurdsoft-challenge/models"
)

func createHandlerFn(s *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			data = &models.Book{}
			book *models.Book
			err  error
		)

		if err = ctx.ShouldBindJSON(data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"parsing json": err.Error()})
			return
		}

		if book, err = s.data.Add(data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, book)
	}
}

func (s *Server) RegisterCreate() {
	s.router.POST(s.config.Root, createHandlerFn(s))
}
