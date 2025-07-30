package models

type ArticleModel struct {
	Model
	Title        string    `gorm:"size:32" json:"title"`
	Abstract     string    `gorm:"size:32" json:"abstract"`
	Content      string    `json:"content"`
	CategoryID   uint      `json:"category_id"`
	TagList      []string  `gorm:"type:longtext;serializer:json" json:"tag_list"`
	Cover        string    `gorm:"size:32" json:"cover"`
	UserID       uint      `json:"user_id"`
	UserModel    UserModel `gorm:"foreignKey:UserID" json:"user_model"`
	LookCount    uint      `json:"look_count"`
	DiggCount    uint      `json:"digg_count"`
	CommentCount string    `json:"commentCount"`
	CollectCount string    `json:"collectCount"`
	OpenComment  bool      `json:"open_comment"`
	Status       int8      `json:"status"` //草稿 审核中 已发布
}
