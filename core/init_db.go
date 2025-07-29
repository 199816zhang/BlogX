package core

import (
	"blogx_server/global"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

func InitDB() *gorm.DB {
	dc := global.Config.DB
	dc1 := global.Config.DB1
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
	if !dc1.Empty() {
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{mysql.Open(dc1.DSN())},//读库	
			Sources:  []gorm.Dialector{mysql.Open(dc.DSN())},//写库
		}))
		if err != nil {
			logrus.Fatalf("读写配置err:%v", err)
		}
	}
	return db
}
