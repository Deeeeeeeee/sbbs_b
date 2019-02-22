package tag

import (
	"log"
	"net/http"
	"sbbs_b/common"
	"sbbs_b/entity"
	"sbbs_b/user"

	"github.com/gin-gonic/gin"
)

// SetupTagAPIRouter 初始化标签 api
func SetupTagAPIRouter(r *gin.RouterGroup) {
	r.GET("/page", tagPage, user.JwtMiddelware())
	r.POST("/publish", publishTag, user.JwtMiddelware())
}

// tagPage 标签分页
func tagPage(c *gin.Context) {
	page := c.GetInt("page")
	size := c.GetInt("size")

	var persist []entity.Tag
	if _, err := common.DBEngine().Limit(size, (page-1)*size).FindAndCount(&persist); err != nil {
		log.Println(err)
		panic(common.HTTP500Error("查询数据库失败"))
	}

	c.JSON(http.StatusOK, persist)
}

// publishTag 发布标签
func publishTag(c *gin.Context) {
	var dto entity.Tag
	common.BindJSON(c, &dto)

	userID, _ := c.Get("userId")
	dto.UserID = userID.(int64)
	dto.ViewCount = 0
	dto.PraiseCount = 0

	var id int64
	var err error
	if id, err = common.DBEngine().InsertOne(dto); err != nil {
		log.Println(err)
		panic(common.HTTP500Error("创建标签失败"))
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
