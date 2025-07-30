package models

import "time"

type UserModel struct {
	Model
	Username       string `gorm:"size:32" json:"username"`
	Nickname       string `gorm:"size:32" json:"nickname"`
	Avatar         string `gorm:"size:256" json:"avatar"`
	Abstract       string `gorm:"size:256" json:"abstract"`
	RegisterSource int8   `json:"register_source"`
	CodeAge        int    `json:"code_age"`
	Password       string `gorm:"size:64" json:"-"`
	Email          string `gorm:"size:256" json:"email"`
	OpenID         string `gorm:"size:64" json:"open_id"` //用户在那个第三方平台上的唯一身份标识
	Role           int8   `json:"role"`
}
type UserConfModel struct {
	UserID             uint       `gorm:"unique" json:"user_id"`
	UserModel          UserModel  `gorm:"foreignKey:UserID" json:"user_model"`
	LikeTags           []string   `gorm:"type:longtext;serializer:json" json:"like_tags"`
	UpdateUsernameDate *time.Time `json:"update_username_date"`
	OpenCollect        bool       `json:"open_collect"`
	OpenFollow         bool       `json:"open_follow"`
	OpenFans           bool       `json:"open_fans"`
	HomeStyleID        uint       `json:"home_style_id"`
}
