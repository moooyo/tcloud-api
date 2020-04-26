package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func practiceRouter(g *gin.RouterGroup) {
	group := g.Group("/practice")
	group.Use(middlewares.AuthRequired)
	group.POST("", controllers.CreatePractice)
	group.GET("", controllers.GetPracticeList)
	group.PATCH("/:id", controllers.PatchPractice)
}