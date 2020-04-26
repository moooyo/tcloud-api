package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func courseRouter(r *gin.RouterGroup) {
	group := r.Group("/course")
	group.Use(middlewares.AuthRequired)
	group.POST("", controllers.CreateCourse)
	group.GET("", controllers.GetCourseList)
	group.GET("/directory", controllers.GetCourseDirectory)
	group.PATCH("/:id", controllers.PatchCourse)
	return
}
