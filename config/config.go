package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// T config
type T struct {
	Database struct {
		DataSourceName string `yaml:"dataSourceName"`
		DriverName     string `yaml:"driverName"`
	}
}

// Config 全局变量
var Config T

// InitConfig 初始化配置文件
func InitConfig(env string) {
	var data []byte
	if env == "local" {
		data, _ = ioutil.ReadFile("config/local.yml")
	} else if env == "dev" {
		data, _ = ioutil.ReadFile("config/dev.yml")
	} else if env == "pro" {
		data, _ = ioutil.ReadFile("config/pro.yml")
	}
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalln("配置文件读取失败: " + err.Error())
	}
}
