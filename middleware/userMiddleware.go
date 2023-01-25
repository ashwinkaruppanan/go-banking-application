package middleware

import (
	"net/http"

	"ashwin.com/go-banking-project/helper"
	"github.com/gin-gonic/gin"
)

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, "Token not found")
			c.Abort()
			return
		}

		claims, msg := helper.ValidateToken(token)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		c.Set("name", claims.Name)
		c.Set("user_id", claims.Uid)

		c.Next()
	}
}
