package routers

import (
	"ashwin.com/go-banking-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	router.Use(middleware.UserMiddleware())
	//router.GET("/api/v1/accounts/:user_id", controller.GetAccount())
}
