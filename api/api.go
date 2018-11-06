package api

import (
	"log"
	"net/http"
	"sbbs_b/dao"

	"github.com/gin-gonic/gin"
)

// 一些公共的方法先放这里

// GetOne 根据 id 获取
func GetOne(entity interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		if has, _ := dao.Orm.ID(id).Get(entity); has == true {
			c.JSON(http.StatusOK, entity)
			return
		}
		c.JSON(http.StatusOK, nil)
	}
}

// Create 新增
func Create(entity interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(entity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
			return
		}
		if _, err := dao.Orm.Insert(entity); err != nil {
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
		if err := c.ShouldBindJSON(entity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
			return
		}
		if has, _ := dao.Orm.ID(c.Params.ByName("id")).Get(persist); has == false {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id 对应的数据不存在"})
			return
		}
		dao.Orm.ID(c.Params.ByName("id")).Update(entity)
		c.JSON(http.StatusOK, nil)
	}
}

// Delete 删除用户信息
func Delete(entity interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		dao.Orm.ID(c.Params.ByName("id")).Delete(entity)
		c.JSON(http.StatusOK, nil)
	}
}
