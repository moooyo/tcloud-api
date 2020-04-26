package controllers

import (
	"fmt"
	"net/http"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"

	"github.com/gin-gonic/gin"
)

type createCourseParams struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description" binding:"required"`
	Files       []uint             `json:"files"`
	StartTime   uint               `json:"startTime" binding:"required"`
	EndTime     uint               `json:"endTime" binding:"required"`
	Tags        []models.TagParams `json:"tags"`
	Class       []uint             `json:"class"`
}

func CreateCourse(c *gin.Context) {
	var params createCourseParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form. ID not set.", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	if user == nil || user.Type == 0 {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.UnauthorizedOperation, "Unauthorized operation", ""))
		return
	}
	var course models.Course
	course.Name = params.Name
	course.Description = params.Description
	course.StartTime = params.StartTime
	course.EndTime = params.EndTime
	course.UID = user.ID
	saveCourse, err := models.CreateCourse(&course, params.Tags, params.Files, params.Class)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", saveCourse))
}

type getCourseListParams struct {
	Limit  uint `form:"limit" binding:"required"`
	Offset uint `form:"offset"`
	Course uint `form:"course"`
}

func GetCourseList(c *gin.Context) {
	var params getCourseListParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	var err error
	var ret interface{}
	if params.Course == 0 {
		ret, err = models.GetCourseList(user, params.Offset, params.Limit)
	} else if user != nil {
		v, e := models.GetCourseResponseByID(params.Course)
		f := false
		for _, x := range v.Class {
			if x.ID == user.Class {
				f = true
			}
		}
		if !f {
			e = fmt.Errorf("%s", "No authorized.")
		}
		err = e
		if err == nil {
			ret = v
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
	}
	return
}

type getCourseDirectoryParams struct {
	Course uint `form:"course" binding:"required"`
	Path   uint `form:"path" binding:"required"`
}

func GetCourseDirectory(c *gin.Context) {
	var params getCourseDirectoryParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error()))
		return
	}
	user := middlewares.GetUserInfo(c)
	ret, err := models.GetCourseDirectoryByCourseAndPath(user, params.Course, params.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
	}
}

type patchCourseClassParams struct {
	Class []uint `json:"class" binding:"required"`
}

func PatchCourse(c *gin.Context) {
	var id idParams
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	var params patchCourseClassParams
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
		ret, err = models.PatchCourseClass(user, id.ID, params.Class)
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
