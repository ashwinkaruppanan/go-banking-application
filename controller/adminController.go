package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"ashwin.com/go-banking-project/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ActivateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newOperation *model.ActivateAccount

		if err := c.ShouldBindJSON(&newOperation); err != nil {
			log.Panic(err)
		}

		account_id, err := primitive.ObjectIDFromHex(newOperation.AccountID)
		if err != nil {
			log.Panic(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		res, dberr := accountCollection.UpdateByID(ctx, account_id, bson.M{"$set": bson.M{
			"account_status": newOperation.Operation,
		}})
		defer cancel()
		if dberr != nil {
			log.Panic(dberr)
		}

		if res.MatchedCount == 0 {
			c.JSON(http.StatusBadRequest, "Invalid account ID")
			return
		}

		c.JSON(http.StatusOK, map[string]string{newOperation.AccountID: newOperation.Operation})

	}
}

func ActivateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		var newOperation *model.ActivateUser

		err := c.ShouldBindJSON(&newOperation)
		if err != nil {
			log.Panic(err)
		}

		userID, err := primitive.ObjectIDFromHex(newOperation.UserID)
		if err != nil {
			log.Panic(err)
		}

		res, derr := userCollection.UpdateByID(ctx, userID, bson.M{"$set": bson.M{
			"user_status": newOperation.Operation,
		}})
		defer cancel()
		if derr != nil {
			c.JSON(http.StatusBadRequest, derr.Error())
			return
		}
		if res.MatchedCount == 0 {
			c.JSON(http.StatusBadRequest, "Invalid user ID")
			return
		}

		c.JSON(http.StatusOK, map[string]int{newOperation.UserID: int(newOperation.Operation)})

	}
}
