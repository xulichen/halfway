package db

import (
	"fmt"
	"log"
	"time"

	apmmysql "go.elastic.co/apm/module/apmgormv2/driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySqlConfig struct {
	DNS string `json:"dns"`
	// Host Port User Password DB build DNS
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`
	// Debug开关
	Debug bool `json:"debug" yaml:"debug" default:"false"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns" yaml:"maxIdleConns" default:"10"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns" yaml:"maxOpenConns" default:"100"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" yaml:"connMaxLifetime" default:"0"`
}

// NewDBWithAPM 初始化 MySQL 链接
func NewDBWithAPM(cf *MySqlConfig) *gorm.DB {
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
	if cf.Debug {
		mdb = mdb.Debug()
	}
	if cf.MaxIdleConns > 0 {
		sqlDB, _ := mdb.DB()
		sqlDB.SetMaxIdleConns(cf.MaxIdleConns)
	}
	if cf.MaxOpenConns > 0 {
		sqlDB, _ := mdb.DB()
		sqlDB.SetMaxOpenConns(cf.MaxOpenConns)
	}
	if cf.ConnMaxLifetime > 0 {
		sqlDB, _ := mdb.DB()
		sqlDB.SetConnMaxLifetime(cf.ConnMaxLifetime)
	}

	log.Println("connected")
	return mdb
}

// NewDB 初始化 MySQL 链接
func NewDB(cf *MySqlConfig) *gorm.DB {
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
	if cf.Debug {
		mdb = mdb.Debug()
	}
	if cf.MaxIdleConns > 0 {
		sqlDB, _ := mdb.DB()
		sqlDB.SetMaxIdleConns(cf.MaxIdleConns)
	}
	if cf.MaxOpenConns > 0 {
		sqlDB, _ := mdb.DB()
		sqlDB.SetMaxOpenConns(cf.MaxOpenConns)
	}
	if cf.ConnMaxLifetime > 0 {
		sqlDB, _ := mdb.DB()
		sqlDB.SetConnMaxLifetime(cf.ConnMaxLifetime)
	}
	log.Println("connected")
	return mdb
}
