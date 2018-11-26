package dao

import (
	"log"
	"os"
	"sbbs_b/config"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	orm  *xorm.Engine
	once sync.Once
)

// Engine 返回 xorm
func Engine() *xorm.Engine {
	lazyinit()
	return orm
}

func lazyinit() {
	once.Do(func() {
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
	err = orm.Sync2(new(User), new(Tag), new(Comment))
	if err != nil {
		log.Fatalln(err)
	}
}
