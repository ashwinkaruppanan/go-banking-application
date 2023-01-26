package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"ashwin.com/go-banking-project/database"
	"ashwin.com/go-banking-project/helper"
	"ashwin.com/go-banking-project/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func AdminMiddleware() gin.HandlerFunc {
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

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		adminID, err := primitive.ObjectIDFromHex(claims.Uid)
		if err != nil {
			log.Panic(err)
		}

		var dbRecord *model.User
		err = userCollection.FindOne(ctx, bson.M{"_id": adminID}).Decode(&dbRecord)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, "Admin id not found")
			c.Abort()
			return
		}

		if dbRecord.UserType != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()

			return
		}
		c.Set("user_id", claims.Uid)
		c.Set("name", claims.Name)
		c.Next()
	}
}
