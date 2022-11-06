package blog

import "blog/global"

//如果考虑到多用户 其实也不用担心
// user->article->tag
type Tag struct {
	global.GM_MODEL
	Name    string    `json:"name"gorm:"comment:标签"`
	Article []Article `gorm:"many2many:blog_article_tag;"`
	Uid     uint      `json:"u_id"` //多用户
}

type Category struct {
	global.GM_MODEL
	Uid     uint      `json:"u_id"` //多用户
	Name    string    `json:"name"gorm:"comment:分类"`
	Article []Article `gorm:"many2many:blog_article_category;"`
}

func (t LikeAndWatch) TableName() string {
	return "blog_count_like_watched"
}
func (t Tag) TableName() string {
	return "blog_tag"
}
func (a Category) TableName() string {
	return "blog_category"
}
