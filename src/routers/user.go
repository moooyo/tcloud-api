package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func usersRouter(g *gin.RouterGroup) {
	group := g.Group("/users")
	group.POST("", controllers.Register)
	group.POST("/confirm", controllers.RegisterConfirm)
	group.GET("/code", controllers.GenerateRegisterCode)
	group.GET("/", middlewares.WithAuthRequired(controllers.GetUserList))
	userRouter(g)
}

func userRouter(g *gin.RouterGroup) {
	group := g.Group("/user")
	group.Use(middlewares.AuthRequired)
	group.PATCH("/:id", controllers.UpdateUserInfo)
}
