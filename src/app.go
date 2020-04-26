package main

import (
	"log"
	"tcloud-api/src/models"
	"tcloud-api/src/routers"
	"tcloud-api/src/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const dev = true

func main() {
	util.InitLogger()
	util.WorkSpaceInit()
	defer util.WorkSpaceClean()
	util.InitMailServer()

	config := util.GetConfig().Web
	r := gin.Default()

	// init
	models.InitDataBase()
	r.Use(cors.Default())

	router := r.Group("/api")
	routers.InitRouters(router)

	err := r.Run(util.FormatUrl(config.Address, config.Port))
	if err != nil {
		log.Fatal(err)
	}
}
