package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"

	"github.com/gin-gonic/gin"
)

func GetTagsList(c *gin.Context) {
	u := middlewares.GetUserInfo(c)
	if u == nil {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.UnauthorizedOperation, "Unauthorized operation.", ""))
		return
	}
	q := c.DefaultQuery("q", "")
	var ret interface{}
	var err error
	var create bool
	if q != "" {
		ret, create, err = models.GetTagsByNameWithCreate(q)
	} else {
		ret, err = models.GetTagsList()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, fmt.Sprint(create), ret))
	}
	return
}

func GetTagsByID(c *gin.Context) {
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
	tags, err := models.GetTagsByID(ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavaliable.", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", tags))
	}
	return
}
