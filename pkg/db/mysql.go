package db

import (
	"fmt"
	"log"

	apmmysql "go.elastic.co/apm/module/apmgormv2/driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/xulichen/halfway/pkg/config"
)

// NewDBWithAPM 初始化 MySQL 链接
func NewDBWithAPM(cf *config.MySqlConfig) *gorm.DB {
	log.Println("connecting MySQL ... ", cf.Host,
		cf.Port,
		cf.DB)
	mdb, err := gorm.Open(apmmysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cf.User, cf.Password, cf.Host, cf.Port, cf.DB)), &gorm.Config{},
	)
	if err != nil {
		panic(err)
	}
	if mdb == nil {
		panic("failed to connect database")
	}
	log.Println("connected")
	return mdb
}

// NewDB 初始化 MySQL 链接
func NewDB(cf *config.MySqlConfig) *gorm.DB {
	log.Println("connecting MySQL ... ", cf.Host,
		cf.Port,
		cf.DB)
	mdb, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cf.User, cf.Password, cf.Host, cf.Port, cf.DB)), &gorm.Config{},
	)
	if err != nil {
		panic(err)
	}
	if mdb == nil {
		panic("failed to connect database")
	}
	log.Println("connected")
	return mdb
}

// SetDebug ...
func SetDebug(db *gorm.DB, on bool) *gorm.DB {
	db = db.Debug()
	return db
}
