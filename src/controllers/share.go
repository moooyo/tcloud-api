package controllers

import (
	"encoding/base64"
	"net/http"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"
	"time"

	"github.com/gin-gonic/gin"
)

type createShareParams struct {
	PathID  []uint `json:"path" binding:"required"`
	Secret  bool   `json:"secret"`
	Name    string `json:"name" binding:"required"`
	Expired int64  `json:"expired"`
}

type createShareResponse struct {
	ID        string `json:"id"`
	ShareName string `json:"shareName"`
	Secret    bool   `json:"secret"`
	Password  string `json:"password"`
	Expired   int64  `json:"expired"`
}

func CreateShare(c *gin.Context) {
	//var params createShareParams
	var params createShareParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Create share error.", err.Error()))
		return
	}

	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusForbidden, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", ""))
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", "User assume error"))
		return
	}

	expired := time.Unix(params.Expired/1000, 0)

	uuid, password, err := models.CreateShare(&user, params.Name, params.PathID, expired, params.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", createShareResponse{
		ShareName: params.Name,
		ID:        uuid,
		Secret:    params.Secret,
		Password:  password,
		Expired:   params.Expired,
	}))
	return
}

type shareinfoParams struct {
	ID string `uri:"id"`
}

type shareResponse struct {
	IsOwner  bool                    `json:"owner"`
	NickName string                  `json:"nickname"`
	ID       uint                    `json:"id"`
	Secret   bool                    `json:"secret"`
	FileList []models.ShareDirectory `json:"fileList"`
	Expired  int64                   `json:"expired"`
	RootID   uint                    `json:"path"`
	ShareID  uint                    `json:"shareID"`
}

func GetShareBaseInfo(c *gin.Context) {
	var params shareinfoParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Get share info error.", err.Error()))
		return
	}
	tmp, err := base64.StdEncoding.DecodeString(params.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Get share info error.", err.Error()))
		return
	}
	var resp shareResponse
	uuid := string(tmp)
	record, err := models.GetShareRecord(uuid)
	timestamp := record.Expired.Unix() * 1000
	resp.Expired = timestamp
	resp.Secret = record.Secret
	resp.ID = record.ShareRootID
	resp.NickName = record.NickName
	resp.IsOwner = false
	resp.RootID = record.ShareRootID
	resp.ShareID = record.ID

	user := middlewares.GetUserInfo(c)
	if user != nil && user.ID == record.UID {
		resp.IsOwner = true
	}

	password := c.DefaultQuery("code", "")
	if password != "" && record.Password != password {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.UnauthorizedOperation, "Password incorrect.", ""))
		return
	}

	if resp.IsOwner || !record.Secret || (record.Secret && password == record.Password) {
		dirs, err := models.GetShareDirectoriesByPreIndex((record.ShareRootID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Get share info error.", err.Error()))
			return
		}
		resp.FileList = dirs
	}

	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", resp))
}

type shareListParams struct {
	Offset uint `form:"offset"`
	Limit  uint `form:"limit" binding:"required"`
	Path   uint `form:"path"`
}

func GetShareList(c *gin.Context) {
	var params shareListParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Get share record error.", err.Error()))
		return
	}

	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusForbidden, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", ""))
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", "User assume error"))
		return
	}

	if params.Path == 0 {
		dir, err := models.GetShareRecordListByUserID(user.ID, params.Offset, params.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
			return
		}
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", dir))
	} else {
		dir, err := models.GetShareDirectoriesByPreIndex(params.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable.", err.Error()))
			return
		}
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", dir))
	}
}
