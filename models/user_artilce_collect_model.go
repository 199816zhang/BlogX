package models

import "time"

type UserArtilceCollectModel struct {
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"user_id"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"user_model"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"article_id"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"article_model"`
	CreatedAt    time.Time    `json:"created_at"`
	CollectID    uint         `gorm:"uniqueIndex:idx_name" json:"collect_id"`
	CollectModel CollectModel `gorm:"foreignKey:CollectID" json:"collect_model"`
}
