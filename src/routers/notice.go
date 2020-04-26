package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func noticeRouter(g *gin.RouterGroup) {
	group := g.Group("notice")
	group.Use(middlewares.AuthRequired)
	group.GET("", controllers.GetNoticeList)
	group.POST("", controllers.CreateNotice)
	group.PATCH("/:id", controllers.PatchNotice)
}
