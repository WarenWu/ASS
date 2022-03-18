package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"ASS/config"
)

// var engine *xorm.Engine

// func InitDatabse() {
// 	dataSourceName := config.MysqlUser + ":" + config.MysqlPassword + "@tcp(" + config.MysqlAddr + ")/" + config.MysqlDb + "?charset=utf8"
// 	e, err := xorm.NewEngine("mysql", dataSourceName)
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	if config.Debug {
// 		f, err := os.Create("./sql.log")
// 		if err != nil {
// 			logrus.Fatalln(err)
// 		}
// 		engine.SetLogger(log.NewSimpleLogger(f))
// 		engine.Logger().SetLevel(log.LOG_DEBUG)
// 		engine.ShowSQL(true)
// 	}
// 	err = e.Sync2(new(StockInfo))
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	err = e.Sync2(new(StockCommonInfo))
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	engine = e
// 	DbEngine()
// }

// func DbEngine() *xorm.Engine {
// 	return engine
// }

var database *gorm.DB

func InitDatabse() {
	dataSourceName := config.MysqlUser + ":" + config.MysqlPassword + "@tcp(" + config.MysqlAddr + ")/" + config.MysqlDb + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		logrus.Fatalln(err)
	}
	db.AutoMigrate(&StockInfo{},&StockCommonInfo{})
	database = db
}

func DbEngine() *gorm.DB {
	return database
}
