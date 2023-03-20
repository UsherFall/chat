package db

import (
	"chat/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var dbMap = map[string]*gorm.DB{}

func init() {
	initDB("gochat")
}

//初始化mysql连接
func initDB(dbName string) {
	var e error
	dbConfig := config.Conf.Db
	dbMap[dbName], e = gorm.Open("mysql", dbConfig.DbBase.Link)
	dbMap[dbName].DB().SetMaxIdleConns(4)  //连接池初始数
	dbMap[dbName].DB().SetMaxOpenConns(20) //连接池最大数
	dbMap[dbName].DB().SetConnMaxLifetime(8 * time.Second)
	if e != nil {
		logrus.Error("connect db fail:%s", e.Error())
	}
}

func GetDb(dbName string) (db *gorm.DB) {
	if db, ok := dbMap[dbName]; ok {
		return db
	} else {
		return nil
	}
}

type DbGoChat struct {
}

func (*DbGoChat) GetDbName() string {
	return "gochat"
}
