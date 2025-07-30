package models

type CommentModel struct {
	Model
	Content        string          `gorm:"size:256" json:"content"`
	UserID         uint            `json:"user_id"`
	UserModel      UserModel       `gorm:"foreignKey:UserID" json:"-"`
	ArticleID      uint            `json:"article_id"`
	ArticleModel   ArticleModel    `gorm:"foreignKey:ArticleID" json:"-"`
	ParentID       *uint           `json:"parent_id"` //父评论
	ParentModel    *CommentModel   `gorm:"foreignKey:ParentID" json:"-"`
	SubCommentList []*CommentModel `gorm:"foreignKey:ParentID" json:"-"`
	RootParentID   *uint           `json:"root_parent_id"`
	DiggCount      uint            `json:"digg_count"`
}
