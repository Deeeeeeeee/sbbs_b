package common

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

// 一些公共的方法先放这里

// BindJSON 绑定 JSON 参数
func BindJSON(c *gin.Context, dto interface{}) {
	if err := c.ShouldBindJSON(dto); err != nil {
		handleValidateErr(c, err, "binding")
	}
}

// BindJSONWithValidate 绑定参数并校验
func BindJSONWithValidate(c *gin.Context, dto interface{}, tagName string) {
	if err := c.ShouldBindJSON(dto); err != nil {
		handleValidateErr(c, err, "binding")
	}
	v, _ := SingletonValidator(tagName)
	if err := v.Struct(dto); err != nil {
		handleValidateErr(c, err, tagName)
	}
}

func handleValidateErr(c *gin.Context, err error, tagName string) {
	var message string
	// 如果是参数校验失败，返回中文提示信息
	if errs, ok := err.(validator.ValidationErrors); ok {
		message = fmt.Sprintf("%v", errs.Translate(GetZHTrans(tagName)))
	} else {
		message = err.Error()
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message:": "参数错误:" + message})
	c.Error(err)
	panic(err)
}

// GetOne 根据 id 获取
func GetOne(entity interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		if has, _ := DBEngine().ID(id).Get(entity); has == true {
			c.JSON(http.StatusOK, entity)
			return
		}
		c.JSON(http.StatusOK, nil)
	}
}

// Create 新增
func Create(entity interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		BindJSON(c, entity)
		if _, err := DBEngine().Insert(entity); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, entity)
	}
}

// Update 更新
func Update(entity interface{}, persist interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		BindJSON(c, entity)
		if has, _ := DBEngine().ID(c.Params.ByName("id")).Get(persist); has == false {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id 对应的数据不存在"})
			return
		}
		DBEngine().ID(c.Params.ByName("id")).Update(entity)
		c.JSON(http.StatusOK, nil)
	}
}

// Delete 删除用户信息
func Delete(entity interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		DBEngine().ID(c.Params.ByName("id")).Delete(entity)
		c.JSON(http.StatusOK, nil)
	}
}
