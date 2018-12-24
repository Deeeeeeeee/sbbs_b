package common

import (
	"log"
	"os"
	"os/user"
	"sbbs_b/comment"
	"sbbs_b/config"
	"sbbs_b/tag"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	orm     *xorm.Engine
	ormOnce sync.Once
)

// DBEngine 返回 xorm
func DBEngine() *xorm.Engine {
	dbLazyinit()
	return orm
}

func dbLazyinit() {
	ormOnce.Do(func() {
		initOrm()
	})
}

func initOrm() {
	database := config.Config.Database
	if database.DriverName == "sqlite3" {
		os.Remove(database.DataSourceName)
	}

	var err error
	orm, err = xorm.NewEngine(database.DriverName, database.DataSourceName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	orm.ShowSQL(true)
	orm.SetMaxIdleConns(20)
	orm.SetMaxOpenConns(50)

	// 同步表结构
	err = orm.Sync2(new(user.User), new(tag.Tag), new(comment.Comment))
	if err != nil {
		log.Fatalln(err)
	}
}
