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
	r.POST("", Create(&dao.User{}))
	r.PUT("/:id", Update(&dao.User{}, &dao.User{}))
	r.DELETE("/:id", Delete(&dao.User{}))
}

// UserPage 分页获取用户信息
func userPage(c *gin.Context) {
	var users []dao.User
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	dao.Orm.Limit(size, (page-1)*size).Find(&users)
	c.JSON(http.StatusOK, users)
}
