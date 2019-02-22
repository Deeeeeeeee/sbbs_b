package service

import (
	"log"
	"sbbs_b/common"
	"sbbs_b/entity"
)

// CheckTagExists 检查是否存在标签
func CheckTagExists(tagID int64) entity.Tag {
	var persist entity.Tag
	if exists, err := common.DBEngine().ID(tagID).Get(&persist); err != nil {
		log.Println(err)
		panic(common.HTTP500Error("执行sql语句错误:" + err.Error()))
	} else if exists == false {
		panic(common.HTTP400Error("对应标签不存在"))
	}
	return persist
}
