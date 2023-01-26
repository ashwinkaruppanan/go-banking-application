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
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var accountCollection *mongo.Collection = database.OpenCollection(database.Client, "account")
var transactionCollection *mongo.Collection = database.OpenCollection(database.Client, "transaction")

func CreateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {

		var newAccount *model.Account

		if err := c.ShouldBindJSON(&newAccount); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		user_id := c.GetString("user_id")

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

func GetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dbDetails *model.Account
		userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
		if err != nil {
			log.Panic(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err = accountCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&dbDetails)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		showDetails := model.ShowAccount{
			AccountType:   dbDetails.AccountType,
			Balance:       dbDetails.Balance,
			AccountStatus: dbDetails.AccountStatus,
		}

		c.JSON(http.StatusOK, showDetails)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
		if err != nil {
			log.Panic(err)
		}
		var userDetails *model.User

		err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&userDetails)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		showUser := model.ShowUser{
			FullName:   userDetails.FullName,
			Email:      userDetails.Email,
			UserStatus: userDetails.UserStatus,
		}

		c.JSON(http.StatusOK, showUser)

	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var dbDetails *model.User
		var newDetails *model.UpdateUser

		if err := c.ShouldBindJSON(&newDetails); err != nil {
			log.Panic(err)
		}

		if validateErr := validator.New().Struct(newDetails); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr})
			return
		}

		UserID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
		if err != nil {
			log.Panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err = userCollection.FindOne(ctx, bson.M{"_id": UserID}).Decode(&dbDetails); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		defer cancel()

		//verifypassword
		isCorrect := VerifyPassword(newDetails.CurrentPassword, dbDetails.Password)
		if !isCorrect {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Password"})
			return
		}

		if newDetails.FullName != "" {
			dbDetails.FullName = newDetails.FullName
		}

		if newDetails.Email != "" {
			dbDetails.Email = newDetails.Email
		}

		if newDetails.NewPassword != "" {
			newHashedPassword := HashPassword(newDetails.NewPassword)
			dbDetails.Password = newHashedPassword
		}

		dbDetails.UpdatedAt = time.Now().Unix()

		filter := bson.M{"_id": UserID}
		update := bson.M{"$set": bson.M{
			"email":      dbDetails.Email,
			"full_name":  dbDetails.FullName,
			"password":   dbDetails.Password,
			"updated_at": dbDetails.UpdatedAt,
		}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if newDetails.NewPassword != "" && newDetails.Email == "" && newDetails.FullName == "" {
			c.JSON(http.StatusOK, "password updated!")
			return
		}
		out := map[string]string{"full_name": dbDetails.FullName, "email": dbDetails.Email}

		c.JSON(http.StatusOK, out)

	}
}

func Transfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserId, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
		if err != nil {
			log.Panic(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		var newTransactionData model.Transaction
		var newTransaction *model.Transfer

		if err := c.ShouldBindJSON(&newTransaction); err != nil {
			log.Panic(err)
		}

		var sender *model.User

		err = userCollection.FindOne(ctx, bson.M{"_id": UserId}).Decode(&sender)
		defer cancel()
		if err != nil {
			log.Panic(err)
		}

		if sender.UserStatus != 1 {
			c.JSON(http.StatusBadRequest, "your not allowed for transaction")
			return
		}

		var senderDetails *model.Account
		var receiverDetails *model.Account
		err = accountCollection.FindOne(ctx, bson.M{"user_id": UserId}).Decode(&senderDetails)
		defer cancel()
		if err != nil {
			log.Panic(err)
		}

		desAccountID, err := primitive.ObjectIDFromHex(newTransaction.DES_Account)
		if err != nil {
			log.Panic(err)
		}
		err = accountCollection.FindOne(ctx, bson.M{"_id": desAccountID}).Decode(&receiverDetails)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid receiver account id"})
			return
		}

		if senderDetails.AccountStatus != "ACTIVE" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "your account is under verfication"})
			return
		}

		if senderDetails.Balance < newTransaction.Amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
			return
		}

		if desAccountID == senderDetails.AccountID {
			c.JSON(http.StatusBadRequest, "self transfer not allowed")
			return
		}

		if receiverDetails.AccountStatus != "ACTIVE" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "receiver's account is under verfication"})
			return
		}

		_, err = accountCollection.UpdateOne(ctx, bson.M{"_id": senderDetails.AccountID}, bson.M{
			"$set": bson.M{
				"balance": senderDetails.Balance - newTransaction.Amount,
			}})
		defer cancel()

		if err != nil {
			log.Panic(err)
		}

		_, err = accountCollection.UpdateOne(ctx, bson.M{"_id": receiverDetails.AccountID},
			bson.M{"$set": bson.M{
				"balance": receiverDetails.Balance + newTransaction.Amount,
			}})
		defer cancel()
		if err != nil {
			log.Panic(err)
		}

		des_id, err := primitive.ObjectIDFromHex(newTransaction.DES_Account)
		if err != nil {
			log.Panic(err)
		}

		newTransactionData.TransactionID = primitive.NewObjectID()
		newTransactionData.UserID = UserId
		newTransactionData.TransactionTime = time.Now().Unix()
		newTransactionData.SRC_Account = senderDetails.AccountID
		newTransactionData.DES_Account = des_id
		newTransactionData.Operation = "DEBIT"
		newTransactionData.Amount = newTransaction.Amount

		_, err = transactionCollection.InsertOne(ctx, newTransactionData)
		defer cancel()
		if err != nil {
			log.Panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"Transaction ID": newTransactionData.TransactionID,
		})

	}
}

func ViewTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		userID := c.GetString("user_id")
		userIDObject, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Panic(err)
		}
		var userDetails *model.Account
		if err := accountCollection.FindOne(ctx, bson.M{"user_id": userIDObject}).Decode(&userDetails); err != nil {
			log.Panic(err)
		}
		defer cancel()

		var debit []*model.Transaction
		var credit []*model.Transaction
		cursor, err := transactionCollection.Find(ctx, bson.M{"src_account": userDetails.AccountID})
		defer cancel()
		if err != nil {
			log.Panic(err)
		}

		if err = cursor.All(ctx, &debit); err != nil {
			log.Panic(err)
		}

		cursor, err = transactionCollection.Find(ctx, bson.M{"des_account": userDetails.AccountID})
		defer cancel()
		if err != nil {
			log.Panic(err)
		}

		if err = cursor.All(ctx, &credit); err != nil {
			log.Panic(err)
		}

		for _, val := range credit {
			val.Operation = "CREDIT"
		}

		debit = append(debit, credit...)
		c.JSON(http.StatusOK, debit)
	}
}
