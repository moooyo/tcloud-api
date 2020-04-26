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

type params struct {
	ID uint `uri:"id"`
}

type fileListParams struct {
	Offset uint `form:"offset"`
	Limit  uint `form:"limit" binding:"required"`
	PathID uint `form:"path" binding:"required"`
	Type   uint `form:"type"`
}

func GetFileList(c *gin.Context) {
	var params fileListParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "invalid form", ""))
		return
	}

	user, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.ContextInfoNotSet, "User context info not set", ""))
		return
	}
	u, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.ContextInfoNotSet, "User assert error", ""))
		return
	}

	var data []models.Directory
	var err error
	// Type === 0 means not set
	if params.Type == 0 {
		data, err = models.SearchFileListByPreIndex(params.PathID, params.Offset, params.Limit)
	} else {
		data, err = models.SearchFileListByType(&u, params.Type, params.Offset, params.Limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", data))
	}
}

type createDirectoryParams struct {
	ID   uint   `json:"path"`
	Name string `json:"name"`
}

func CreateDirectory(c *gin.Context) {
	var params createDirectoryParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
		return
	}
	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", ""))
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Assume error.", ""))
		return
	}
	directory, err := models.CreateDirectoryByPathID(params.ID, user.ID, params.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", ""))
		return
	}

	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", directory))
	return
}

type changeFileNameParams struct {
	Name string `json:"name" binding:"required"`
}
type changeFileNameIDParams struct {
	ID uint `uri:"id" binding:"required"`
}

func ChangeFileName(c *gin.Context) {
	var params changeFileNameParams
	var id changeFileNameIDParams
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
		return
	}
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
		return
	}

	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", ""))
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Assume error.", ""))
		return
	}

	if err := models.ChangeDirectoryName(&user, id.ID, params.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Operation error.", err.Error))
		return
	}

	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ""))
}

type deleteFileParams struct {
	IDstr string `uri:"id"`
}

func DeleteFile(c *gin.Context) {
	var params deleteFileParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
		return
	}
	completelyStr := c.DefaultQuery("completely", "false")
	completely, err := strconv.ParseBool(completelyStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
		return
	}

	IDstr := strings.Split(params.IDstr, ",")
	var id []uint
	for _, str := range IDstr {
		num, err := strconv.Atoi(str)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
			return
		}
		id = append(id, uint(num))
	}

	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", ""))
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Assume error.", ""))
		return
	}

	if !completely {
		if err := models.DeleteDirectories(id, user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Delete file error.", err.Error()))
			return
		}
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ""))
	} else {
		ret, err := models.DeleteFileCompletely(&user, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Delete file completely error.", err.Error()))
			return
		}
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
	}

}

type trashListPrams struct {
	Limit  uint `form:"limit"`
	Offset uint `form:"offset"`
}

func GetTrashList(c *gin.Context) {
	var params trashListPrams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
		return
	}

	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", ""))
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Assume error.", ""))
		return
	}

	ret, err := models.GetDeletedFileList(&user, params.Offset, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavailable", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ret))
}

type patchParams struct {
	IDstr string `uri:"id"`
}

func PatchFile(c *gin.Context) {
	var params deleteFileParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
		return
	}
	IDstr := strings.Split(params.IDstr, ",")
	var id []uint
	for _, str := range IDstr {
		num, err := strconv.Atoi(str)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error))
			return
		}
		id = append(id, uint(num))
	}

	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", ""))
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Assume error.", ""))
		return
	}

	operation := c.DefaultQuery("op", "unset")
	switch operation {
	case "restore":
		ids, err := models.RestoreTrash(&user, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavaliable", ids))
			return
		} else {
			c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ids))
			return
		}
	case "unset":
	default:
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "op unset", ""))
		return
	}
}
