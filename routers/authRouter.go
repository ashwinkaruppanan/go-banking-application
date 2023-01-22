package routers

import (
	"ashwin.com/go-banking-project/controller"
	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.Engine) {
	router.POST("/signup", controller.Signup())
	router.POST("/login", controller.Login())
	router.DELETE("/logout", controller.Logout())
}
