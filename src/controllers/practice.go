package controllers

import (
	"net/http"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"

	"github.com/gin-gonic/gin"
)

type createPracticeParams struct {
	OJ uint `json:"oj" binding:"required"`
	ID uint `json:"id" binding:"required"`
}

func CreatePractice(c *gin.Context) {
	var params createPracticeParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	resp, err := models.CreatePractice(user, params.OJ, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", resp))
	}
	return
}

type practiceList struct {
	Limit  uint `form:"limit" binding:"required"`
	Offset uint `form:"offset"`
}

func GetPracticeList(c *gin.Context) {
	var params practiceList
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	var err error
	var ret []models.PracticeResponse
	ret, err = models.GetPracticeList(user, params.Offset, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
	}
	return
}

type patchPracticeClassParams struct {
	Class []uint `json:"class" binding:"required"`
}
type idParams struct {
	ID uint `uri:"id" binding:"required"`
}

func PatchPractice(c *gin.Context) {
	var id idParams
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	var params patchPracticeClassParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	op := c.DefaultQuery("op", "")
	var err error
	var ret interface{}
	switch op {
	case "class":
		ret, err = models.PatchPracticeClass(user, id.ID, params.Class)
	default:
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, " Invalid op.", ""))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
	}
}
