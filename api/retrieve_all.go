package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func retrieveAllHandlerFn(s *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		books := s.data.GetAll()
		ctx.JSON(http.StatusOK, books)
	}
}

func (s *Server) RegisterRetrieveAll() {
	s.router.GET(s.config.Root, retrieveAllHandlerFn(s))
}
