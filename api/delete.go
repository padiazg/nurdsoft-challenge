package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterDelete() {
	s.router.DELETE(s.config.Root+"/:id", func(ctx *gin.Context) {
		var (
			ids = ctx.Param("id")
			err error
			id  int64
		)

		id, err = strconv.ParseInt(ids, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"parsing id to int:": err.Error()})
			return
		}

		err = s.data.Delete(int32(id))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"deleting:": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"delete": fmt.Sprintf("succesfully deleted record %d", id)})
	})
}
