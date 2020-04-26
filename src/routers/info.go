package routers

import (
	"github.com/gin-gonic/gin"
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"
)

func infoRouter(r *gin.RouterGroup) {
	router := r.Group("/info")
	router.Use(middlewares.AuthRequired)
	router.GET("/user", controllers.GetSessionUserInfo)
}
