package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func shareRouter(g *gin.RouterGroup) {
	group := g.Group("/share")
	group.POST("", middlewares.WithAuthRequired(controllers.CreateShare))
	group.GET("/:id", controllers.GetShareBaseInfo)
	group.GET("", middlewares.WithAuthRequired(controllers.GetShareList))
}
