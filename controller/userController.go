package controller

// import (
// 	"context"
// 	"time"

// 	"ashwin.com/go-banking-project/database"
// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

// func GetAccount() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

// 		userId := c.Param("user_id")

// 		err := userCollection.FindOne(ctx, bson.M{"us"})

// 	}
// }
