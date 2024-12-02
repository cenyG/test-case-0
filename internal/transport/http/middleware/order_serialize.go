package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JSONDeserializerMiddleware(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(model); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("jsonData", model)
		c.Next()
	}
}
