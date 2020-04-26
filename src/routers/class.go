package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func classRouters(g *gin.RouterGroup) {
	group := g.Group("/class")
	group.GET("", controllers.GetClassList)
	group.POST("", middlewares.WithAuthRequired(controllers.CreateClassInfo))
	group.PATCH("/:id", middlewares.WithAuthRequired(controllers.UpdateClassInfo))
}
