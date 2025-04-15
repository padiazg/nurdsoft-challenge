package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRetrieveAll() {
	s.router.GET(s.config.Root, func(ctx *gin.Context) {
		books := s.data.GetAll()
		ctx.JSON(http.StatusOK, books)
	})
}
