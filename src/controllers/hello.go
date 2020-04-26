package controllers

import (
	"net/http"
	"tcloud-api/src/util"

	"github.com/gin-gonic/gin"
)

func HelloWorld(c *gin.Context) {
	_, err := util.GetProblem(2, 1001)
	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, util.TCLOUD_API_VERSION, err))
	return
}
