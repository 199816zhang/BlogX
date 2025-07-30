package models

type CategoryModel struct {
	Model
	Title     string    `gorm:"size:128" json:"title"`
	UserID    uint      `json:"user_id"`
	UserModel UserModel `gorm:"foreignKey:UserID" json:"user_model"`
}
