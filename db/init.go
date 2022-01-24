package db

import (
	"os"

	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
	"xorm.io/xorm/log"

	"ASS/config"
)

var engine *xorm.Engine

func IntDatabse() {
	dataSourceName := config.MysqlUser + ":" + config.MysqlPassword + "@tcp(" + config.MysqlAddr + ")/" + config.MysqlDb_default + "?charset=utf8"
	e, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		logrus.Fatalln(err)
	}
	if config.Debug {
		f, err := os.Create("./sql.log")
		if err != nil {
			logrus.Fatalln(err)
		}
		engine.SetLogger(log.NewSimpleLogger(f))
		engine.Logger().SetLevel(log.LOG_DEBUG)
		engine.ShowSQL(true)
	}
	err = e.Sync2(new(StockInfo))
	if err != nil {
		logrus.Fatalln(err)
	}
	engine = e
	DbEngine()
}

func DbEngine() *xorm.Engine {
	return engine
}
