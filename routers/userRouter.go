package routers

import (
	"ashwin.com/go-banking-project/controller"
	"ashwin.com/go-banking-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	router.Use(middleware.UserMiddleware())
	router.POST("/api/v1/user/account", controller.CreateAccount())
}
