package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tcloud-api/src/models"
	"tcloud-api/src/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

const SessionKey = "SESSION_ID"
const ContextUserKey = "CONTEXT_USER"

func AuthRequiredWithReturn(c *gin.Context) bool {
	sessionID, err := c.Cookie(SessionKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			util.FormatResponse(util.CookieNotSet, "please login first.", ""))
		return false
	}

	client := util.GetSessionRedisClient()
	defer client.Close()
	u, err := client.Get(sessionID).Result()
	if err == redis.Nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			util.FormatResponse(util.CookieNotSet, "login status has expired", ""))
		return false
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			util.FormatResponse(util.RedisUnavailable, "server unavailable", ""))
		return false
	}
	//refresh expire
	err = client.Expire(sessionID, 30*time.Minute).Err()
	if err != nil {
		util.ERROR("%s", err.Error())
	}
	var user models.User
	err = json.Unmarshal([]byte(u), &user)
	c.Set(ContextUserKey, user)
	c.Next()
	return true
}

func AuthRequired(c *gin.Context) {
	AuthRequiredWithReturn(c)
}

func WithAuthRequired(handle gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ret := AuthRequiredWithReturn(c)
		if ret {
			handle(c)
		}
	}
}

func GetUserInfo(c *gin.Context) *models.User {
	sessionID, err := c.Cookie(SessionKey)
	if err != nil {
		return nil
	}

	client := util.GetSessionRedisClient()
	u, err := client.Get(sessionID).Result()
	if err != nil {
		return nil
	}
	//refresh expire
	err = client.Expire(sessionID, 30*time.Minute).Err()
	if err != nil {
		util.ERROR("%s", err.Error())
	}
	var user models.User
	err = json.Unmarshal([]byte(u), &user)
	if err != nil {
		return nil
	}
	return &user
}

func GetContextUser(c *gin.Context) (*models.User, error) {
	u, ok := c.Get(ContextUserKey)
	if !ok {
		return nil, fmt.Errorf("%s", "Please login first.")
	}
	user, ok := u.(models.User)
	if !ok {
		return nil, fmt.Errorf("%s", "Assume error.")
	}
	return &user, nil
}
