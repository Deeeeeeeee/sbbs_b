package dao

import (
	"log"
	"os"
	"sbbs_b/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

// Orm xorm engine
var Orm *xorm.Engine

// InitOrm 初始化 orm
func InitOrm() {
	database := config.Config.Database
	if database.DriverName == "sqlite3" {
		os.Remove(database.DataSourceName)
	}

	var err error
	Orm, err = xorm.NewEngine(database.DriverName, database.DataSourceName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	Orm.ShowSQL(true)
	Orm.SetMaxIdleConns(20)
	Orm.SetMaxOpenConns(50)

	// 同步表结构
	err = Orm.Sync2(new(User), new(Tag), new(Comment))
	if err != nil {
		log.Fatalln(err)
	}
}
