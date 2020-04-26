package controllers

import (
	"net/http"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"
	"time"

	"github.com/gin-gonic/gin"
)

type registerForm struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Nickname   string `json:"nickname" binding:"required"`
	ClassName  uint   `json:"class"`
	InviteCode string `json:"class_code"`
	Type       int    `json:"type"`
	TypeCode   string `json:"type_code"`
}

func Register(c *gin.Context) {
	var form registerForm
	// check register form
	if err := c.ShouldBind(&form); err != nil {
		util.WARN("%s", err.Error())
		c.JSON(http.StatusOK, util.FormatResponse(util.InvalidForm, "invalid form", ""))
		return
	}
	if form.ClassName != 0 && form.InviteCode == "" {
		c.JSON(http.StatusOK, util.FormatResponse(util.InvalidForm, "invalid form", "invite code is missing"))
		return
	}
	if form.Type != 0 && form.TypeCode == "" {
		c.JSON(http.StatusOK, util.FormatResponse(util.InvalidForm, "invalid form", "type code is missing"))
		return
	}
	user := models.User{
		Nickname: form.Nickname,
		Email:    form.Email,
		Password: form.Password,
		Class:    form.ClassName,
		Type:     form.Type,
		Status:   models.UserStatusCreateByRegister,
	}
	err := models.InsertUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, util.FormatResponse(util.DatabaseInsertFailed, "create user failed", err))
		return
	}

	// insert success
	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "insert success", ""))

	// send email && insert register code to redis
	registerCode := util.GenerateRegisterCode()
	client := util.GetRegisterRedisClient()
	defer client.Close()
	err = client.Set(user.Email, registerCode, 5*time.Minute).Err()
	if err != nil {
		util.ERROR("%e", err)
		return
	}
	util.SendRegisterCode(form.Email, registerCode)

	return
}

type confirmForm struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname"`
}

func RegisterConfirm(c *gin.Context) {
	var form confirmForm
	if err := c.ShouldBind(&form); err != nil {
		util.WARN("%s", err.Error())
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "invalid form", ""))
		return
	}

	redis := util.GetRegisterRedisClient()
	defer redis.Close()
	code, err := redis.Get(form.Email).Result()
	if err != nil {
		util.ERROR("%s", err.Error())
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.RedisUnavailable, "redis unavailable", ""))
		return
	}
	if code != form.Code {
		c.JSON(http.StatusOK, util.FormatResponse(util.RegisterCodeIncorrect, "register code is incorrect", ""))
		return
	}

	user, err := models.SearchUserByEmail(form.Email)
	if err != nil {
		util.ERROR("%s", err.Error())
		c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "database unavailable", ""))
		return
	}

	if user.Status == models.UserStatusOK ||
		user.Status == models.UserStatusBlocked ||
		user.Status == models.UserStatusDelete ||
		user.Status == models.UserStatusCreateByLogin && form.Nickname == "" {
		c.JSON(http.StatusOK, util.FormatResponse(util.UnauthorizedOperation, "illegal operation", ""))
		return
	}
	if user.Status == models.UserStatusCreateByLogin {
		user.Nickname = form.Nickname
	}
	user.Status = models.UserStatusOK
	err = models.UpdateUser(user)
	if err != nil {
		util.ERROR("%s", err.Error())
		c.JSON(http.StatusOK, util.FormatResponse(util.DataBaseUnavailable, "database unavailable", ""))
		return
	}

	redis = util.GetSessionRedisClient()
	defer redis.Close()
	uuid := util.GenerateUUID()
	err = redis.Set(uuid, user, time.Duration(util.GetConfig().Redis.Expired)*time.Minute).Err()
	if err != nil {
		c.JSON(util.StatusOK, util.FormatResponse(util.RedisUnavailable, err.Error(), ""))
		return
	}
	c.SetCookie(middlewares.SessionKey, uuid, 3600,
		"/", util.GetConfig().Web.Domain, false, false)

	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ""))
	return
}

func GenerateRegisterCode(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusNotFound, util.FormatResponse(util.InvalidForm, "invalid form", ""))
		return
	}
	code := util.GenerateRegisterCode()
	redis := util.GetRegisterRedisClient()
	defer redis.Close()
	err := redis.Set(email, code, 5*time.Minute).Err()
	if err != nil {
		util.ERROR("%s", err.Error())
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.RedisUnavailable, "redis unavailable", ""))
	}
	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ""))
}
