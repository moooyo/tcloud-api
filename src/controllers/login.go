package controllers

import (
	"net/http"
	"tcloud-api/src/middlewares"
	"tcloud-api/src/models"
	"tcloud-api/src/util"
	"time"

	"github.com/gin-gonic/gin"
)

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Type     int    `form:"type"`
}

func Login(c *gin.Context) {
	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		util.INFO("%s", err.Error())
		c.JSON(http.StatusOK, util.FormatResponse(util.InvalidForm, "invalid form", ""))
		return
	}
	user, loginSuccess := models.SearchUserByLoginForm(form.Email, form.Password, form.Type)

	if !loginSuccess {
		// default value is zero
		// that means user not exist
		if user.ID == 0 {
			user.Status = models.UserStatusCreateByLogin
			user.Password = form.Password
			user.Email = form.Email
			user.Type = form.Type
			err := models.InsertUser(user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, util.FormatResponse(util.DataBaseUnavailable, "create user failed", err))
				return
			}
			c.JSON(http.StatusOK, util.FormatResponse(util.CreateUserByLogin, "", ""))
			code := util.GenerateRegisterCode()
			redis := util.GetRegisterRedisClient()
			defer redis.Close()
			err = redis.Set(form.Email, code, 5*time.Minute).Err()
			if err != nil {
				util.ERROR("%s", err.Error())
			} else {
				util.SendRegisterCode(form.Email, code)
			}
			return
		} else {
			c.JSON(http.StatusOK, util.FormatResponse(util.PasswordIncorrect, "username or password incorrect", ""))
			return
		}
	}

	uuid := util.GenerateUUID()
	redis := util.GetSessionRedisClient()
	defer redis.Close()
	err := redis.Set(uuid, user, time.Duration(util.GetConfig().Redis.Expired)*time.Minute).Err()
	if err != nil {
		c.JSON(util.StatusOK, util.FormatResponse(util.RedisUnavailable, err.Error(), ""))
		return
	}

	c.SetCookie(middlewares.SessionKey, uuid, 3600,
		"/", util.GetConfig().Web.Domain, false, false)

	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ""))

	return
}

func LogOut(c *gin.Context) {
	sessionID, err := c.Cookie(middlewares.SessionKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized,
			util.FormatResponse(util.CookieNotSet, "please login first.", ""))
		return
	}

	client := util.GetSessionRedisClient()
	defer client.Close()
	_, err = client.Del(sessionID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.FormatResponse(util.RedisUnavailable, "Log out error.", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.FormatResponse(util.StatusOK, "", ""))
	return
}
