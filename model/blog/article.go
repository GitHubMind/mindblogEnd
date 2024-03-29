package blog

import (
	"blog/global"
	"html/template"
)

type Article struct {
	global.GM_MODEL
	Title                string        `json:"title" `
	Desc                 string        `json:"desc"`
	Content              template.HTML `json:"content"gorm:"type:longtext"`
	CoverImageUrl        string        `json:"cover_image_url"`
	Uid                  uint          `json:"u_id"`
	State                uint8         `json:"status"gorm:"comment:发布状态,1发布，2未发布;default 2"`
	Tag                  []Tag         `json:"tag"gorm:"many2many:blog_article_tag;comment:标签"`
	Category             []Category    `json:"category"gorm:"many2many:blog_article_category;"`
	LikeAndWatchs        []LikeAndWatch
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
