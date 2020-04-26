package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"

	"github.com/gin-gonic/gin"
)

func GetSessionUserInfo(c *gin.Context) {
	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Context info not set", ""))
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "User assert error", ""))
		return
	}

	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", user))
	return
}

type userListParams struct {
	Offset uint `form:"offset"`
	Limit  uint `form:"limit" binding:"required"`
}

func GetUserList(c *gin.Context) {
	var params userListParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	user, err := middlewares.GetContextUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", err.Error()))
		return
	}

	data, err := models.GetUserList(user, params.Offset, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavaliable.", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", data))
	}

	return
}

type idStrParams struct {
	ID string `uri:"id"`
}

type userInfoParams struct {
	NickName string `json:"nickname" ShouldBind:"required"`
	Email    string `json:"email" ShouldBind:"required"`
	Password string `json:"password"`
	Class    uint   `json:"class"`
}

type classParams struct {
	Class uint `json:"class" ShouldBind:"required"`
}

func UpdateUserInfo(c *gin.Context) {
	var params idStrParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	var idStrArray = strings.Split(params.ID, ",")
	var ids []uint
	for _, v := range idStrArray {
		val, err := strconv.ParseUint(v, 10, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
			return
		}
		ids = append(ids, uint(val))
	}
	op := c.DefaultQuery("op", "unset")
	users, err := models.GetUsersByID(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.DataBaseUnavailable, "Datavase unavailable", err.Error()))
		return
	}

	switch op {
	case "full":
		var fullParams userInfoParams
		if err := c.ShouldBindJSON(&fullParams); err != nil {
			c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
			return
		}
		users[0].Class = fullParams.Class
		users[0].Email = fullParams.Email
		users[0].Nickname = fullParams.NickName
		if fullParams.Password != "" {
			users[0].Password = fullParams.Password
		}
		ret, err := models.UpdateUsers(users)
		code := http.StatusOK
		icode := util.StatusOK
		imessage := ""
		if err != nil {
			code = http.StatusInternalServerError
			icode = util.DataBaseUnavailable
			imessage = "Database unavailable"
		}
		c.JSON(code, util.FormatResponse(icode, imessage, ret))
		break
	case "ban":
		for i := range users {
			users[i].Status = 4
		}
		ret, err := models.UpdateUsers(users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
		} else {
			c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
		}
		break
	case "class":
		var params classParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
			return
		}
		for i := range users {
			users[i].Class = params.Class
		}
		ret, err := models.UpdateUsers(users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
		} else {
			c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
		}
	}

}
