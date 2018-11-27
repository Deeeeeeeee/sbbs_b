package api

import (
	"net/http"
	"sbbs_b/dao"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupUserAPIRouter 初始化 user api router
func SetupUserAPIRouter(r *gin.RouterGroup) {
	r.GET("/:id", GetOne(&dao.User{}))
	r.GET("", userPage)
	r.POST("", userRegistered)
	r.PUT("/:id", Update(&dao.User{}, &dao.User{}))
	r.DELETE("/:id", Delete(&dao.User{}))
}

// UserPage 分页获取用户信息
func userPage(c *gin.Context) {
	var users []dao.User
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	dao.Engine().Limit(size, (page-1)*size).Find(&users)
	c.JSON(http.StatusOK, users)
}

// userRegistered 用户注册
func userRegistered(c *gin.Context) {
	var dto dao.User
	BindJSONWithValidate(c, &dto, "user_register")

	// 根据邮箱查询，如果有则抛出异常，否则新增用户
	if has, _ := dao.Engine().Table("user").Where("email = ?", dto.Email).Exist(); has {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email 已经存在"})
		return
	}
	// 密码加密

	id, _ := dao.Engine().Insert(dto)
	c.JSON(http.StatusOK, gin.H{"id": id})
}
