package http_service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/search/utils/logger"
	"golang.org/x/sys/unix"
	"net"
	"syscall"
)

func RunHttpService() {

	var lc = net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var opErr error
			if err := c.Control(func(fd uintptr) {
				opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			}); err != nil {
				return err
			}
			return opErr
		},
		KeepAlive: 0,
	}

	for i := 0; i < 5; i++ {
		go func(i int) {
			gin.SetMode(gin.ReleaseMode)
			s := Routers(fmt.Sprintf("%d", i))
			lis, err := lc.Listen(context.Background(), "tcp", "0.0.0.0:14444")
			if err != nil {
				panic("启动 http api 失败, err =  " + err.Error())
			}
			logger.Info("启动 Http API , ID:", i)
			err = s.RunListener(lis)
			if err != nil {
				panic("启动 http api 失败, err =  " + err.Error())
			}
		}(i)
	}

}
