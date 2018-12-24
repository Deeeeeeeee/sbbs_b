package main

import (
	"flag"
	"fmt"
	"sbbs_b/common"
	"sbbs_b/config"
	"sbbs_b/user"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func setupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	fmt.Println(config.Config.Log.File)
	// f, _ := os.Create(config.Config.Log.File)
	// gin.DefaultWriter = io.MultiWriter(f)
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery(), common.ExceptionMiddelware())

	// user api
	user.SetupUserAPIRouter(r.Group("/user"))

	return r
}

func setupConfig() {
	// 默认 local
	env := flag.String("env", "local", "环境配置参数")
	flag.Parse()

	// validator
	binding.Validator = common.NewValidator()

	// 初始化配置信息
	config.InitConfig(*env)
}

func main() {
	setupConfig()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
