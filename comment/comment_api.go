package comment

import (
	"log"
	"net/http"
	"sbbs_b/common"
	"sbbs_b/entity"
	"sbbs_b/service"
	"sbbs_b/user"

	"github.com/gin-gonic/gin"
)

// SetupTagAPIRouter 初始化标签 api
func SetupTagAPIRouter(r *gin.RouterGroup) {
	r.GET("/page", commentPage, user.JwtMiddelware())
	r.POST("/comment", comment, user.JwtMiddelware())
}

// commentPage 评论分页
func commentPage(c *gin.Context) {
	page := c.GetInt("page")
	size := c.GetInt("size")
	tagID := c.GetInt64("tagId")

	var persist []entity.Comment
	if _, err := common.DBEngine().Where("tag_id", tagID).
		Limit(size, (page-1)*size).FindAndCount(&persist); err != nil {
		log.Println(err)
		panic(common.HTTP500Error("查询数据库失败"))
	}

	c.JSON(http.StatusOK, persist)
}

// comment 评论
func comment(c *gin.Context) {
	var dto entity.Comment
	common.BindJSON(c, &dto)

	// 判断 tag 是否存在
	service.CheckTagExists(dto.TagID)

	userID, _ := c.Get("userId")
	dto.UserID = userID.(int64)

	var id int64
	var err error
	if id, err = common.DBEngine().InsertOne(dto); err != nil {
		log.Println(err)
		panic(common.HTTP500Error("创建评论失败" + err.Error()))
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
