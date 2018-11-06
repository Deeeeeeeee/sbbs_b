package dao

import "time"

// User describes a user
type User struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}
