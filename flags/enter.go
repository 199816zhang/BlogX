package flags

import (
	"blogx_server/global"
	"blogx_server/models"
	"flag"
	"github.com/sirupsen/logrus"
)

type Options struct {
	File    string
	DB      bool
	Version bool
}

var FlagOptions = new(Options)

func Parse() {
	flag.StringVar(&FlagOptions.File, "f", "settings.yaml", "配置文件")
	flag.BoolVar(&FlagOptions.DB, "db", false, "数据库迁移")
	flag.BoolVar(&FlagOptions.Version, "version", false, "版本信息")
	flag.Parse()
}
func Run() {
	if FlagOptions.DB {
		err := global.DB.AutoMigrate(&models.UserModel{},
			&models.UserConfModel{},
			&models.ArticleModel{},
			&models.ArticleDiggModel{},
			&models.CollectModel{},
			&models.CategoryModel{},
			&models.UserArtilceCollectModel{},
			&models.ImageModel{},
			&models.UserTopArticleModel{},
			&models.UserArticleLookHistoryModel{},
			&models.CommentModel{},
			&models.BannerModel{},
			&models.LogModel{},
			&models.GlobalNotificationModel{},
			&models.UserLoginModel{})
		if err != nil {
			logrus.Fatalf("数据库迁移失败%v", err)
		}
		logrus.Infof("数据库迁移成功")
	}
}
