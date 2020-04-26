package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func tagRouter(c *gin.RouterGroup) {
	group := c.Group("/tag")
	group.Use(middlewares.AuthRequired)
	group.GET("", controllers.GetTagsList)
	group.GET("/:id", controllers.GetTagsByID)
}
