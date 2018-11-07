package dao

import "time"

// User describes a user
type User struct {
	ID          int64     `json:"id"`
	AccountName string    `xorm:"varchar(64) comment '昵称'" json:"accountName"`
	Nickename   string    `xorm:"varchar(64) comment '昵称'" json:"nickname"`
	Password    string    `xorm:"varchar(64) commnet '登录密码'" json:"password"`
	Created     time.Time `xorm:"created" json:"created"`
	Updated     time.Time `xorm:"updated" json:"updated"`
}

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

// Comment describes a comment 评论
type Comment struct {
	ID      int64     `json:"id"`
	UserID  int64     `xorm:"int(20) comment '用户id'" json:"userId"`
	TagID   int64     `xorm:"int(20) comment '标签id'" json:"tagId"`
	Content string    `xorm:"varchar(255) comment '内容'" json:"content"`
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}
