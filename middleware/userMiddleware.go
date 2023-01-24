package middleware

import (
	"net/http"
	"time"

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

		c.SetCookie("user_id", claims.Uid, int(time.Now().Add(5*time.Minute).Unix()), "/", "localhost", false, true)
		c.Next()
	}
}
