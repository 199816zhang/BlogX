package core

import (
	"blogx_server/global"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// InitDB 初始化数据库连接，并根据配置设置读写分离。
// 该函数在程序启动时被调用，返回一个配置好的 *gorm.DB 实例。
func InitDB() *gorm.DB {
	// 从全局配置中获取主数据库（写库）的配置信息。
	dc := global.Config.DB
	// 从全局配置中获取从数据库（读库）的配置信息。
	dc1 := global.Config.DB1
	// 使用GORM连接到主数据库。
	// 所有的写操作（Create, Update, Delete）将默认发送到这个数据库。
	// 同时，它也作为读写分离插件的基础连接。
	db, err := gorm.Open(mysql.Open(dc.DSN()), &gorm.Config{
		// 在执行数据库迁移（AutoMigrate）时，禁用外键约束。
		// 这可以防止在创建或修改表结构时因外键依赖关系而导致失败。
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// 如果数据库连接失败，这是一个致命错误，程序无法继续运行。
		// 记录致命日志并终止程序。
		logrus.Fatalf("数据库连接失败,err:%v", err)
	}
	// 获取底层的 *sql.DB 对象，以便进行更底层的连接池配置。
	sqlDB, err := db.DB()
	if err != nil {
		// 如果获取 *sql.DB 对象失败，同样是致命错误。
		logrus.Fatalf("数据库连接失败,err:%v", err)
	}
	// --- 连接池配置 ---
	// 设置连接池中的最大空闲连接数。
	// 保持一定数量的空闲连接可以避免因频繁创建和销毁连接而带来的性能开销。
	sqlDB.SetMaxIdleConns(10)
	// 设置数据库的最大打开连接数。
	// 这可以防止应用程序无限制地创建连接，从而耗尽数据库资源。
	sqlDB.SetMaxOpenConns(100)
	// 设置连接可被复用的最大时间。
	// 超过这个时间的连接将被关闭并重新创建，这有助于避免因网络问题或数据库重启导致的陈旧无效连接。
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("数据库连接成功")
	// --- 读写分离配置 ---
	// 检查配置文件中是否提供了从数据库（读库）的信息。
	if !dc1.Empty() {
		// 如果提供了读库配置，则启用GORM的dbresolver插件来实现读写分离。
		err = db.Use(dbresolver.Register(dbresolver.Config{
			// Replicas: 配置所有的读数据库实例（从库）。
			// GORM会将读操作（如Find, First, Scan）负载均衡地路由到这些从库。
			Replicas: []gorm.Dialector{mysql.Open(dc1.DSN())}, //读库
			// Sources: 配置所有的写数据库实例（主库）。
			// GORM会将写操作（如Create, Update, Delete）路由到这些主库。
			Sources: []gorm.Dialector{mysql.Open(dc.DSN())}, //写库
		}))
		if err != nil {
			// 如果在配置读写分离插件时出错，这也是一个致命错误。
			logrus.Fatalf("读写配置err:%v", err)
		}
	}
	// 返回配置完成的 *gorm.DB 对象。
	// 如果配置了读写分离，这个db对象会自动处理读写请求的路由。
	// 如果没有配置，它就只是一个连接到主库的普通数据库对象。
	return db
}
