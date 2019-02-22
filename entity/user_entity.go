package entity

import "time"

// User describes a user
type User struct {
	ID        int64     `json:"id" user_login:"required"`
	Nickename string    `xorm:"varchar(64)" json:"nickname" user_register:"required"`
	Email     string    `xorm:"varchar(255)" json:"email" user_register:"required"`
	Password  string    `xorm:"varchar(64)" json:"password" user_register:"required" user_login:"required"`
	Salt      string    `xorm:"varchar(128)" json:"-"`
	Created   time.Time `xorm:"created" json:"created"`
	Updated   time.Time `xorm:"updated" json:"updated"`
}
