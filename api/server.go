package api

import (
	"github.com/gin-gonic/gin"
	"github.com/padiazg/nurdsoft-challenge/internals"
	"github.com/padiazg/nurdsoft-challenge/models"
)

type Server struct {
	router *gin.Engine
	config *models.Config
	data   *internals.BookList
}

func NewServer(config *models.Config, data *internals.BookList) *Server {
	if config.Port == 0 {
		config.Port = 8000
	}

	if config.Root == "" {
		config.Root = "/books"
	}

	if data == nil {
		data = internals.NewBookList()
	}

	res := &Server{
		config: config,
		router: gin.Default(),
		data:   data,
	}

	res.RegisterPing()
	res.RegisterCreate()
	res.RegisterRetrieveAll()
	res.RegisterRetrieveOne()
	res.RegisterUpdate()
	res.RegisterDelete()

	return res
}

func (s *Server) Run() {
	s.router.Run()
}
