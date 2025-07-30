package models

import "blogx_server/models/enum"

type LogModel struct {
	Model
	LogType   enum.LogType   `json:"log_type"`
	Title     string         `gorm:"size:64" json:"title"`
	Content   string         `json:"content"`
	Level     int8           `json:"level"`
	UserID    uint           `json:"user_id"`
	UserModel UserModel      `gorm:"foreignKey:UserID" json:"-"`
	IP        string         `gorm:"size:64" json:"ip"`
	Addr      string         `gorm:"size:64" json:"addr"`
	IsRead    bool           `json:"is_read"`
	LogStatus bool           `json:"log_status"`
	UserName  string         `gorm:"size:64" json:"user_name"`
	Pwd       string         `gorm:"size:64" json:"pwd"`
	LoginType enum.LoginType `json:"login_type"`
}
