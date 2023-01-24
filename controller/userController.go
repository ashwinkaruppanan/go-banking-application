package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"ashwin.com/go-banking-project/database"
	"ashwin.com/go-banking-project/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var accountCollection *mongo.Collection = database.OpenCollection(database.Client, "account")

func CreateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {

		var newAccount *model.Account

		if err := c.ShouldBindJSON(&newAccount); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		user_id, _ := c.Cookie("user_id")

		id, _ := primitive.ObjectIDFromHex(user_id)
		count, err := accountCollection.CountDocuments(ctx, bson.M{"user_id": id})
		defer cancel()
		if err != nil {
			log.Panic(err)
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not able to create account"})
			return
		}
		fmt.Println(count)

		newAccount.AccountID = primitive.NewObjectID()
		newAccount.AccountStatus = "INACTIVE"

		newAccount.UserID = id
		newAccount.CreatedAt = time.Now().Unix()
		newAccount.UpdatedAt = time.Now().Unix()

		_, err = accountCollection.InsertOne(ctx, newAccount)
		if err != nil {
			log.Panic(err)
		}

		c.JSON(http.StatusOK, newAccount)
	}
}
