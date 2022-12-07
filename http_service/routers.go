package http_service

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/search/utils/logger"
	"net"
	"net/http"
	"strings"
	"time"
)

var Router *gin.Engine

func Routers() *gin.Engine {

	Router = gin.Default()

	V1()

	return Router
}

// HttpMiddleware http中间件
func HttpMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t1 := time.Now().UnixNano()
		ip := GetIP(ctx.Request)
		ctx.Set("tum", t1)
		ctx.Next()
		t2 := time.Now().UnixNano()
		logger.Infof("[HTTP] %v | %v | %vum | %vms", ip, ctx.Request.URL.Path, t2-t1, float64(t2-t1)/1e6)
	}
}

func GetIP(r *http.Request) (ip string) {
	for _, ip := range strings.Split(r.Header.Get("X-Forward-For"), ",") {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	if ip = r.Header.Get("X-Real-IP"); net.ParseIP(ip) != nil {
		return ip
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return "0.0.0.0"
}
