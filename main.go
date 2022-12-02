package main

import (
	"github.com/mangenotwork/search/api"
	"github.com/mangenotwork/search/conf"
	"github.com/mangenotwork/search/http_service"
	"github.com/mangenotwork/search/utils"
)

func InitPath() {
	utils.Mkdir(conf.Conf.DataPath)
}

func main() {
	// 读取配置文件
	conf.InitConf()

	// 初始化缓存
	api.InitCache()

	// 连接集群

	// 启动http api
	http_service.RunHttpService()

	// 启动http manage web

	// 启动tcp service

	// 启动grpc service

	// 启动定时任务

	// 启动检查服务

	select {}

}
