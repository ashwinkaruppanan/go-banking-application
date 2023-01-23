package middleware

import (
	"fmt"
	"net/http"

	"ashwin.com/go-banking-project/helper"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Your not authorized")
			c.Abort()
			return
		}

		if cookieValue == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, msg := helper.ValidateToken(cookieValue)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		fmt.Println(claims)
	}
}
