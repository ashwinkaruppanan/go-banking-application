package routers

import (
	"ashwin.com/go-banking-project/controller"
	"ashwin.com/go-banking-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	router.Use(middleware.UserMiddleware())
	router.POST("/api/v1/user/account", controller.CreateAccount())
	router.GET("/api/v1/user/account", controller.GetAccount())
	router.GET("/api/v1/user/", controller.GetUser())
	router.PUT("/api/v1/user/", controller.UpdateUser())
	router.POST("/api/v1/user/transfer/", controller.Transfer())
	router.GET("/api/v1/user/transaction/", controller.ViewTransaction())
}
