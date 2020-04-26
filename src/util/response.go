package util

import "github.com/gin-gonic/gin"

func FormatResponse(code int, message string, data interface{}) (r gin.H) {
	r = gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}
	return
}
