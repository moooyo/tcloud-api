package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"
	"time"
)

type FileUpload struct {
	ACK   uint
	Count uint
	Path  string
	Size  uint
}

func (user *FileUpload) MarshalBinary() (data []byte, err error) {
	return json.Marshal(user)
}
func (user *FileUpload) UnMarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, user)
}

const createFileKey = "CREATE_FILE_AND_SET_UUID"

type uploadFileParams struct {
	UUID     string `json:"uuid" binding:"required"`
	SliceID  uint   `json:"slice_id" binding:"required"`
	Count    uint   `json:"count" binding:"required"`
	Length   uint   `json:"length" binding:"required"`
	File     string `json:"file" binding:"required"`
	PathID   uint   `json:"path_id" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

var directoryPath = util.GetConfig().File.Path

func UploadSingleFile(c *gin.Context) {
	var params uploadFileParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Invalid form", err.Error()))
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

	redisClient := util.GetUploadRedisClient()
	defer redisClient.Close()
	// create file
	if params.UUID == createFileKey {
		fileName := util.GenerateUUID()
		data := FileUpload{
			ACK:   0,
			Count: params.Count,
			Path:  filepath.Join(directoryPath, fileName),
			Size:  0,
		}
		redisClient.Set(fileName, &data, 5*time.Minute)
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", fileName))
		return
	}

	d, err := redisClient.Get(params.UUID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.RedisUnavailable, "", err.Error()))
		return
	}
	var fileUploadACK FileUpload
	err = json.Unmarshal([]byte(d), &fileUploadACK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.RedisUnavailable, "Json unmarshal error", err.Error()))
		return
	}

	if fileUploadACK.ACK != params.SliceID-1 {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Repeat package", fmt.Sprintf("%d %d", fileUploadACK.ACK, params.SliceID)))
		return
	}

	filePath := filepath.Join(directoryPath, params.UUID)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.FileSystemUnavailable, "Upload file error.", err.Error()))
		return
	}
	defer f.Close()

	data, err := base64.StdEncoding.DecodeString(params.File)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.InvalidForm, "Decode error", err.Error()))
		return
	}

	length := len(data)
	if uint(length) != params.Length {
		c.JSON(http.StatusBadRequest, util.FormatResponse(util.InvalidForm, "Upload size not equal to server received", fmt.Sprintf("%d %d", length, params.Length)))
		return
	}

	n, err := f.Write(data)
	if err != nil || n != int(params.Length) {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.WriteFileError, "Upload file err.", fmt.Sprintf("Slice length is not equal to saved. %d != %d", params.Length, n)))
		err = os.Remove(filePath)
		if err != nil {
			util.ERROR("%s", err)
		}
		redisClient.Del(params.UUID)
		return
	}

	// upload success
	if params.SliceID == params.Count {
		redisClient.Del(params.UUID)
		meta := models.FileMeta{
			RealName: params.UUID,
			Size:     fileUploadACK.Size + params.Length,
		}
		err = models.InsertFileMeta(&meta)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DatabaseInsertFailed, "Create file error", err.Error()))
			return
		}
		// insert directory meta
		directory := models.Directory{
			UID:         user.ID,
			Name:        params.FileName,
			IsDirectory: false,
			PreIndex:    params.PathID,
			MetaID:      meta.ID,
			Type:        models.FileName2FileType(params.FileName),
		}
		if !directory.IsDirectory {
			directory.Size = meta.Size
		}
		err = models.InsertDirectory(&directory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DatabaseInsertFailed, "Insert directory or file error.", err.Error()))
			return
		}
	} else {
		fileUploadACK.ACK = fileUploadACK.ACK + 1
		fileUploadACK.Size += params.Length
		redisClient.Set(params.UUID, &fileUploadACK, 5*time.Minute)
	}
	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", params.SliceID))
	return
}
