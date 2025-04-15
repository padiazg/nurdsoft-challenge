package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/padiazg/nurdsoft-challenge/models"
)

func (s *Server) RegisterCreate() {
	s.router.POST(s.config.Root, func(ctx *gin.Context) {
		var (
			data = &models.Book{}
			id   int32
			err  error
		)

		if err = ctx.ShouldBindJSON(data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"parsing json": err.Error()})
			return
		}

		if id, err = s.data.Add(data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"ID": id})
	})
}
