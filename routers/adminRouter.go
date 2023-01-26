package routers

import (
	"ashwin.com/go-banking-project/controller"
	"ashwin.com/go-banking-project/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRouter(c *gin.Engine) {
	c.Use(middleware.AdminMiddleware())
	c.POST("/api/v1/admin/activate/account/", controller.ActivateAccount())
	c.POST("/api/v1/admin/activate/user/", controller.ActivateUser())
}
