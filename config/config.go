package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Config 全局变量
var Config T

// T config
type T struct {
	Application string
	Database    struct {
		DataSourceName string `yaml:"dataSourceName"`
		DriverName     string `yaml:"driverName"`
	}
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	Log struct {
		File string
	}
}

// InitConfig 初始化配置文件
func InitConfig(env string) {
	var data []byte
	var configFile = ""
	switch env {
	case "local":
		configFile = "config/local.yml"
	case "dev":
		configFile = "config/dev.yml"
	case "prod":
		configFile = "config/prod.yml"
	}
	data, _ = ioutil.ReadFile(configFile)
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalln("配置文件读取失败: " + err.Error())
	}
}
