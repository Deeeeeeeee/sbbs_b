package entity

import "time"

// Tag describes a tag 标签
type Tag struct {
	ID          int64     `json:"id"`
	UserID      int64     `xorm:"int(20) comment '用户id'" json:"userId"`
	Content     string    `xorm:"varchar(255) comment '内容'" json:"content"`
	ViewCount   int64     `xorm:"int(20) comment '浏览数量'" json:"viewCount"`
	PraiseCount int64     `xorm:"int(20) comment '点赞数量'" json:"PraiseCount"`
	IsTop       bool      `xorm:"bit default 0 comment '是否置顶'" json:"isTop"`
	Created     time.Time `xorm:"created" json:"created"`
	Updated     time.Time `xorm:"updated" json:"updated"`
}
