package common

import (
	"log"
	"os"
	"sbbs_b/config"
	"sbbs_b/entity"
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
	err = orm.Sync2(new(entity.User), new(entity.Tag), new(entity.Comment))
	if err != nil {
		log.Fatalln(err)
	}
}
