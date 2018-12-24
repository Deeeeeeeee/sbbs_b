package comment

import "time"

// Comment describes a comment 评论
type Comment struct {
	ID      int64     `json:"id"`
	UserID  int64     `xorm:"int(20) comment '用户id'" json:"userId"`
	TagID   int64     `xorm:"int(20) comment '标签id'" json:"tagId"`
	Content string    `xorm:"varchar(255) comment '内容'" json:"content"`
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}
