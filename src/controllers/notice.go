package controllers

import (
	"net/http"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"

	"github.com/gin-gonic/gin"
)

type createNoticeParams struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Level       uint   `json:"level" binding:"required"`
}

func CreateNotice(c *gin.Context) {
	var params createNoticeParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	resp, err := models.CreateNotice(user, params.Title, params.Description, params.Level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", resp))
	}
	return
}

type noticeListParams struct {
	Limit  uint `form:"limit" binding:"required"`
	Offset uint `form:"offset"`
}

func GetNoticeList(c *gin.Context) {
	var params noticeListParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	var err error
	var ret []models.NoticeResponse
	ret, err = models.GetNoticeList(user, params.Offset, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
	}
	return
}

type patchNoticeClassParams struct {
	Class []uint `json:"class" binding:"required"`
}

func PatchNotice(c *gin.Context) {
	var id idParams
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	var params patchNoticeClassParams
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
		ret, err = models.PatchNoticeClass(user, id.ID, params.Class)
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
