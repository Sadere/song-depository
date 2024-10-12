package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Checks whether Content-Type header is set to application/json
func CheckJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "please set content-type to json",
			})
			return
		}

		c.Next()
	}
}

// Sets Content-Type header to application/json
func JSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		c.Next()
	}
}
