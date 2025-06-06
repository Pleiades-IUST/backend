package ginutil

import (
	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "USER_ID_KEY"
)

func GetUserID(ctx *gin.Context) int64 {
	return ctx.Value(UserIDKey).(int64)
}
