package models

type CollectModel struct {
	Model
	Title        string    `gorm:"size:32" json:"title"`
	Abstract     string    `gorm:"size:256" json:"abstract"`
	Cover        string    `gorm:"size:256" json:"cover"`
	ArticleCount int64     `json:"article_count"`
	UserID       uint      `json:"user_id"`
	UserModel    UserModel `gorm:"foreignKey:UserID" json:"user_model"`
}
