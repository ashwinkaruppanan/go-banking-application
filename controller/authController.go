package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"ashwin.com/go-banking-project/database"
	"ashwin.com/go-banking-project/helper"
	"ashwin.com/go-banking-project/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var tokenCollection *mongo.Collection = database.OpenCollection(database.Client, "token")
var accountCollection *mongo.Collection = database.OpenCollection(database.Client, "account")

func HashPassword(password string) (hashPassword string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(loginPassword string, dbPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(loginPassword))
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {

		var newUser model.CreateAccount
		var newUserDetails model.User
		var newAccountDetails model.Account

		if err := c.BindJSON(&newUser); err != nil {
			log.Panic(err)
		}

		validateErr := validator.New().Struct(newUser)
		if validateErr != nil {
			log.Panic(validateErr)
		}

		newUserDetails.UserID = primitive.NewObjectID()
		newUserDetails.Email = newUser.User.Email
		newUserDetails.Password = HashPassword(newUser.User.Password)
		newUserDetails.FullName = newUser.User.FullName
		newUserDetails.UserType = newUser.User.UserType
		newUserDetails.UserStatus = 0
		newUserDetails.CreatedAt = time.Now().Unix()
		newUserDetails.UpdatedAt = time.Now().Unix()

		newAccountDetails.AccountID = primitive.NewObjectID()
		newAccountDetails.UserID = newUserDetails.UserID
		newAccountDetails.AccountType = newUser.Account.AccountType
		newAccountDetails.Balance = newUser.Account.Balance
		newAccountDetails.AccountStatus = "INACTIVE"
		newAccountDetails.CreatedAt = time.Now().Unix()
		newAccountDetails.CreatedAt = time.Now().Unix()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": newUserDetails.Email})
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()

		if count > 0 {
			c.JSON(http.StatusConflict, "Email already registered with some other account")
			return
		}

		_, err = userCollection.InsertOne(ctx, newUserDetails)
		if err != nil {
			log.Panic(err)
		}
		defer cancel()

		_, err = accountCollection.InsertOne(ctx, newAccountDetails)
		defer cancel()
		if err != nil {
			log.Panic(err)
		}

		c.JSON(http.StatusOK, "ok")

		// 	var newUser *model.User

		// 	if err := c.BindJSON(&newUser); err != nil {
		// 		log.Panic(err)
		// 	}

		// 	validateErr := validator.New().Struct(newUser)
		// 	if validateErr != nil {
		// 		c.JSON(http.StatusBadRequest, validateErr)
		// 		return
		// 	}

		// 	newUser.UserID = primitive.NewObjectID()
		// 	newUser.UserStatus = 0
		// 	newUser.CreatedAt = time.Now().Unix()
		// 	newUser.UpdatedAt = time.Now().Unix()
		// 	newUser.Password = HashPassword(newUser.Password)
		// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// 	count, err := userCollection.CountDocuments(ctx, bson.M{"email": newUser.Email})
		// 	defer cancel()
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	if count > 0 {
		// 		c.JSON(http.StatusConflict, "Email already registered with some other account")
		// 		return
		// 	}

		// 	_, err = userCollection.InsertOne(ctx, newUser)
		// 	defer cancel()

		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	c.JSON(http.StatusOK, newUser)
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginDetails *model.User
		var dbDetails *model.User
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

		if err := c.BindJSON(&loginDetails); err != nil {
			log.Panic(err)
		}

		err := userCollection.FindOne(ctx, bson.M{"email": loginDetails.Email}).Decode(&dbDetails)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is incorrect"})
			log.Panic(err)
		}

		isPasswordCorrect := VerifyPassword(loginDetails.Password, dbDetails.Password)

		if !isPasswordCorrect {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token := helper.GenerateToken(dbDetails.FullName, dbDetails.UserID.Hex(), dbDetails.UserType)
		timeNow := time.Now()
		tokenDB := &model.Token{
			UserID:    dbDetails.UserID,
			Token:     token,
			CreatedAt: timeNow.Unix(),
		}
		_, err = tokenCollection.InsertOne(ctx, tokenDB)
		defer cancel()
		if err != nil {
			log.Fatal(err)
		}
		c.SetCookie("token", token, int(timeNow.Add(24*time.Hour).Unix()), "/", "localhost", false, true)

	}

}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		cxt, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		token, err := c.Request.Cookie("token")
		if err != nil {
			log.Fatal(err)
		}

		tokenCollection.DeleteOne(cxt, bson.M{"token": token.Value})
		defer cancel()

		c.SetCookie("token", "", -1, "/", "", false, true)

		c.JSON(http.StatusOK, "Logged out successfully")
	}
}
