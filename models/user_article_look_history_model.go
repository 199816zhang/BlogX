package models

type UserArticleLookHistoryModel struct {
	Model
	UserID       uint         `json:"user_id"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"user_model"`
	ArticleID    uint         `json:"article_id"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"article_model"`
}
