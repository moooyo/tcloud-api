package controllers

import (
	"net/http"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"

	"github.com/gin-gonic/gin"
)

func GetClassList(c *gin.Context) {
	user := middlewares.GetUserInfo(c)
	list, err := models.GetClassInfoList(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "DataBase unavailable", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", list))
	}

	return
}

type createParams struct {
	Name string `json:"name" binding:"required"`
}

func CreateClassInfo(c *gin.Context) {
	var params createParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	info, err := models.CreateClassInfo(user, params.Name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.UnauthorizedOperation, "Unauthorized operation.", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", info))
	}
	return
}

type updateParams struct {
	ID uint `uri:"id" binding:"required"`
}

func UpdateClassInfo(c *gin.Context) {
	var params updateParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form. ID not set.", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	op := c.DefaultQuery("op", "unset")
	info, err := models.GetClassInfoByID(user, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
		return
	}
	switch op {
	case "code":
		info.Code = util.GenerateCaptcha(models.ClassCodeLength)
		info, err = models.UpdateClassInfo(user, info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
		} else {
			c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", info))
		}
		break
	case "name":
		name := c.DefaultQuery("name", "unset")
		if name == "unset" {
			c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form. Name not set.", ""))
			break
		}
		info.Name = name
		info, err = models.UpdateClassInfo(user, info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
		} else {
			c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", info))
		}
		break
	case "unset":
	default:
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form. OP not set", ""))
		break
	}

	return
}
