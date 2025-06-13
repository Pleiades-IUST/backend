package drive

import "github.com/gin-gonic/gin"

func Register(e *gin.Engine) {
	r := e.Group("drive")

	r.POST("", CreateDrive)
	r.GET("all", FetchAllDrives)
	r.POST("signals", FetchSignals)
}
