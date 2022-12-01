package http_service

import (
	"github.com/gin-gonic/gin"
)

func RunHttpService() {
	go func() {
		gin.SetMode(gin.DebugMode)
		s := Routers()
		s.Run(":14444")
	}()
}
