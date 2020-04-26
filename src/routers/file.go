package routers

import (
	"tcloud-api/src/controllers"
	"tcloud-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func filesRouter(g *gin.RouterGroup) {
	group := g.Group("/files")
	group.Use(middlewares.AuthRequired)
	group.GET("", controllers.GetFileList)
	group.GET("/trash", controllers.GetTrashList)

	group.POST("", controllers.UploadSingleFile)
	group.POST("/download", controllers.FilesDownload)
	group.POST("/directory", controllers.CreateDirectory)
}

func fileRouter(g *gin.RouterGroup) {

	group := g.Group("/file")
	group.Use(middlewares.AuthRequired)
	group.POST("/:id/name", controllers.ChangeFileName)
	group.DELETE("/:id", controllers.DeleteFile)
	group.GET("/:id", controllers.FileDownload)
	group.PATCH("/:id", controllers.PatchFile)
}
