package core

import (
	"blogx_server/global"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitDB() *gorm.DB {
	dc := global.Config.DB
	db, err := gorm.Open(mysql.Open(dc.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败,err:%v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("数据库连接失败,err:%v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("数据库连接成功")
	return db
}
