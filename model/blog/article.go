package blog

import (
	"blog/global"
	"html/template"
)

type Article struct {
	global.GM_MODEL
	Title                string         `json:"title" `
	Desc                 string         `json:"desc"`
	Content              template.HTML  `json:"content"gorm:"type:longtext"`
	CoverImageUrl        string         `json:"cover_image_url"`
	Uid                  uint           `json:"u_id"`
	State                uint8          `json:"state"gorm:"comment:发布状态"`
	Tag                  []Tag          `json:"tag"gorm:"many2many:blog_article_tag;comment:标签"`
	LikeAndWatcks        []LikeAndWatck `json:"like_watched"gorm:"many2many:blog_article_l_w;"`
	Category             []Category     `json:"category"gorm:"many2many:blog_article_category;"`
	ArticleContentBackUp ArticleContentBackUp
}

func (a Article) TableName() string {
	return "blog_article"
}

type ArticleContentBackUp struct {
	Content   template.HTML `json:"content"gorm:"type:longtext"`
	ArticleID uint
}

func (a ArticleContentBackUp) TableName() string {
	return "blog_ArticleContentBackUp"
}
