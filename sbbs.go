package main

import (
	"flag"
	"sbbs_b/api"
	"sbbs_b/config"
	"sbbs_b/dao"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	r := gin.Default()

	// user api
	api.SetupUserAPIRouter(r.Group("/user"))

	return r
}

func setupConfig() {
	// 默认 local
	env := flag.String("env", "local", "环境配置参数")
	flag.Parse()

	// 初始化配置信息
	config.InitConfig(*env)
	// 初始化 orm
	dao.InitOrm()
}

func main() {
	setupConfig()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
