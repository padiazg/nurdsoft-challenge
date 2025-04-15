package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/padiazg/nurdsoft-challenge/internals"
)

func deleteHandlerFn(s *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			ids        = ctx.Param("id")
			err        error
			id         int64
			statusCode int
		)

		id, err = strconv.ParseInt(ids, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"parsing id to int:": err.Error()})
			return
		}

		err = s.data.Delete(int32(id))
		if err != nil {
			if errors.As(err, &internals.ErrorNotFound{}) {
				statusCode = http.StatusNotFound
			} else {
				statusCode = http.StatusBadRequest
			}

			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"delete": fmt.Sprintf("succesfully deleted record %d", id)})
	}
}

func (s *Server) RegisterDelete() {
	s.router.DELETE(s.config.Root+"/:id", deleteHandlerFn(s))
}
