package conf

import (
	"github.com/mangenotwork/search/utils"
	"github.com/mangenotwork/search/utils/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Conf *config

type config struct {
	HttpService *HttpService `yaml:"http_service"`
	DataPath    string       `yaml:"data_path"`
}

type HttpService struct {
	Prod string `yaml:"prod"`
}

// InitConf 读取yaml文件
// 获取配置
func InitConf() {
	appConfigPath := "configs.yaml"
	if !utils.FileExists(appConfigPath) {
		panic("【启动失败】 未找到配置文件!")
	}
	logger.Info("[启动]读取配置文件:", appConfigPath)
	//读取yaml文件到缓存中
	config, err := ioutil.ReadFile(appConfigPath)
	if err != nil {
		panic("【启动失败】读取配置文件" + err.Error())
	}
	err = yaml.Unmarshal(config, &Conf)
	if err != nil {
		panic("【启动失败】读取配置文件" + err.Error())
	}

	utils.Mkdir(Conf.DataPath + "/doc/")
	utils.Mkdir(Conf.DataPath + "/index/")

}
