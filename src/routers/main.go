package routers

import (
	"tcloud-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouters(g *gin.RouterGroup) {
	//CORS config
	g.GET("/hello", controllers.HelloWorld)
	sessionRouters(g)
	usersRouter(g)
	infoRouter(g)
	filesRouter(g)
	fileRouter(g)
	shareRouter(g)
	classRouters(g)
	tagRouter(g)
	courseRouter(g)
	practiceRouter(g)
	noticeRouter(g)
}
