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

// var tokenCollection *mongo.Collection = database.OpenCollection(database.Client, "token")

func HashPassword(password string) (hashPassword string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(loginPassword string, dbPassword string) bool {
	out := true
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(loginPassword))
	if err != nil {
		// log.Panic(err)
		out = false
	}
	return out
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {

		var newUser *model.User

		if err := c.BindJSON(&newUser); err != nil {
			log.Panic(err)
		}

		validateErr := validator.New().Struct(newUser)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, validateErr)
			return
		}

		newUser.UserID = primitive.NewObjectID()
		if newUser.UserType == "ADMIN" {
			newUser.UserStatus = 1
		} else {
			newUser.UserStatus = 0
		}
		newUser.CreatedAt = time.Now().Unix()
		newUser.UpdatedAt = time.Now().Unix()
		newUser.Password = HashPassword(newUser.Password)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": newUser.Email})
		defer cancel()
		if err != nil {
			log.Fatal(err)
		}

		if count > 0 {
			c.JSON(http.StatusConflict, "Email already registered with some other account")
			return
		}

		_, err = userCollection.InsertOne(ctx, newUser)
		defer cancel()

		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, newUser)
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		}

		token := helper.GenerateToken(dbDetails.FullName, dbDetails.UserID.Hex(), dbDetails.UserType)

		// tokenDB := &model.Token{
		// 	UserID:    dbDetails.UserID,
		// 	Token:     token,
		// 	CreatedAt: timeNow.Unix(),
		// }
		// _, err = tokenCollection.InsertOne(ctx, tokenDB)
		// defer cancel()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		c.SetCookie("token", token, 60*60*3, "/", "localhost", false, true)
		c.JSON(http.StatusOK, "Login Successfull")
	}

}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// cxt, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		// token, err := c.Request.Cookie("token")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// tokenCollection.DeleteOne(cxt, bson.M{"token": token.Value})
		// defer cancel()

		c.SetCookie("token", "", -1, "/", "", false, true)

		c.JSON(http.StatusOK, "Logged out successfully")
	}
}
