package auth

import (
	authmiddleware "github.com/Pleiades-IUST/backend/webservice/middleware/auth"
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(e *gin.Engine) {
	r := e.Group("auth")

	r.POST("signup", signup)
	r.POST("login", login)
	r.GET("protected", authmiddleware.AuthenticateUser, protected)
}
