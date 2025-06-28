package drive

import (
	"github.com/Pleiades-IUST/backend/webservice/middleware/auth"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	r := e.Group("drive")

	r.POST("", auth.AuthenticateUser, CreateDrive)
	r.GET("all", auth.AuthenticateUser, FetchAllDrives)
	r.POST("signals", auth.AuthenticateUser, FetchSignals)
	r.GET("csv", GetCSV)
}
