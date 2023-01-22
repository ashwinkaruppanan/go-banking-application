package main

import (
	"log"
	"os"

	"ashwin.com/go-banking-project/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	routers.AuthRouter(router)
	routers.UserRouter(router)
	routers.AdminRouter(router)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	router.Run(":" + port)

}
