package controllers

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"

	"github.com/gin-gonic/gin"
)

type downloadParams struct {
	FileID []uint `json:"files" binding:"required"`
}

func walkDir(c *gin.Context, d *models.Directory, path string) error {
	if d.IsDirectory {
		nextPath := filepath.Join(path, d.Name)
		util.CreateDirectoryWhenNotExist(nextPath)
		list, err := models.SearchFileListByPreIndexWithoutLimit(d.ID)
		if err != nil {
			return err
		}
		for _, dir := range list {
			err := walkDir(c, &dir, nextPath)
			if err != nil {
				return err
			}
		}
	} else {
		config := util.GetConfig().File
		meta, err := models.SearchFileMetaByID(d.MetaID)
		if err != nil {
			util.ERROR("%s", err.Error())
			return err
		}
		// create hard link
		metaPath := filepath.Join(config.Path, meta.RealName)
		targetPath := filepath.Join(path, d.Name)
		cmd := exec.Command("ln", metaPath, targetPath)
		err = cmd.Run()
		if err != nil {
			util.ERROR("%s", err.Error())
			return err
		}
	}
	return nil
}

func FilesDownload(c *gin.Context) {
	var params downloadParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
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
	basePath := util.GetConfig().File.Path
	tmpPath := filepath.Join(basePath, "tmp")
	var templateFilePath = filepath.Join(tmpPath, "download")

	uuid := util.GenerateUUID()
	dirPath := filepath.Join(templateFilePath, uuid)
	util.CreateDirectoryWhenNotExist(dirPath)

	defer util.RemoveFileOrDirectoryWithoutError(dirPath)

	for _, dir := range params.FileID {
		d, err := models.SearchDirectoryByID(dir)
		if err != nil {
			util.ERROR("%s", err.Error())
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Files damaged", err.Error()))
			return
		}
		if d.UID != user.ID {
			c.JSON(http.StatusUnauthorized, util.FormatResponse(util.UnauthorizedOperation, "Illegal operation", ""))
			return
		}
		err = walkDir(c, d, dirPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.FileDamaged, "Download files error!", err.Error()))
			return
		}
	}

	//pack all file to zip
	packPath := filepath.Join(templateFilePath, uuid+".tar.gz")
	defer util.RemoveFileOrDirectoryWithoutError(packPath)

	cmd := exec.Command("tar", "-zcvf", packPath, "-C", dirPath, ".")
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.FileDamaged, "Download files error!", err.Error()))
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "pack.tar.gz"))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(packPath)
	return
}

type downloadFileParams struct {
	ID uint `uri:"id"`
}
type shareDownloadParams struct {
	Share uint `form:"share"`
}

type shareCourseParams struct {
	Course uint `form:"course"`
}

func FileDownload(c *gin.Context) {
	var params downloadFileParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
		return
	}
	var file *models.FileMeta
	var err error
	var name string
	u, ok := c.Get(middlewares.ContextUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.FormatResponse(util.ContextInfoNotSet, "Please login first.", ""))
		return
	}
	op := c.DefaultQuery("op", "unset")
	switch op {
	case "unset":
		user, ok := u.(models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.ContextInfoNotSet, "Assume error.", ""))
			return
		}

		dir, err := models.SearchDirectoryByID(params.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavaliable.", err.Error))
			return
		}

		if dir.UID != user.ID || dir.IsDirectory {
			c.JSON(http.StatusUnauthorized, util.FormatResponse(util.UnauthorizedOperation, "Unauthorized operation.", ""))
			return
		}

		file, err = models.SearchFileMetaByID(dir.MetaID)
		name = dir.Name
		break
	case "share":
		{
			var shareParams shareDownloadParams
			if err := c.ShouldBind(&shareParams); err != nil {
				c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
				return
			}
			file, name, err = models.GetMetaIDByDirectoryAndShare(params.ID, shareParams.Share)
		}
		break
	case "course":
		{
			u := middlewares.GetUserInfo(c)
			var courseParams shareCourseParams
			if err := c.ShouldBind(&courseParams); err != nil {
				c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form.", err.Error()))
				return
			}
			file, name, err = models.GetShareCourseFileMeta(u, courseParams.Course, params.ID)
			break
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "Database unavaliable.", err.Error()))
		return
	}
	config := util.GetConfig().File
	path := filepath.Join(config.Path, file.RealName)

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(path)
}
