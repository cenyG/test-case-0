package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type fail struct {
	Error string `json:"error" example:"message"`
}

type success struct {
	Status string `json:"status"`
	ID     uint64 `json:"id"`
}

func failResponse(c *gin.Context, code int, msg string) {
	slog.Error(fmt.Sprintf("[failResponse] error: %s", msg))
	c.AbortWithStatusJSON(code, fail{msg})
}

func successResponse(c *gin.Context, id uint64) {
	c.JSON(200, success{
		Status: "OK",
		ID:     id,
	})
}
