package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func sessionRouters(r *gin.RouterGroup) {
	router := r.Group("/session")
	router.POST("", controllers.Login)
	router.DELETE("", middlewares.WithAuthRequired(controllers.LogOut))
}
