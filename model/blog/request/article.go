package request

import (
	"blog/model/blog"
	"blog/model/commond/request"
	"html/template"
)

type ArticleRequset struct {
	blog.Article
}
type ArticleRequsetID struct {
	ID uint `json:"id"`
}
type ArticleRequsetIDMul struct {
	ID []ArticleRequsetID `json:"IDS"`
}
type ArticleContent struct {
	Content template.HTML `json:"content"`
	ID      uint          `json:"id"`
}

// 标题查询  状态查询
type ArticleSearchRequset struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	State       uint8  `json:"status"`
	ArticleTags []uint `json:"Tag"`
	Category    []uint `json:"category""`
	request.PageInfo
}
