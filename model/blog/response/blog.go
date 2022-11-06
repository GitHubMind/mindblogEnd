package response

import "blog/model/blog"

type BlogInfoResponse struct {
	ArticleCount  int64  `json:"article_count";`
	TagCount      int64  `json:"tag_count";`
	CategoryCount int64  `json:"category_count";`
	GitHubAddress string `json:"github_address"`
	Uid           uint   `json:"u_id"`
	HeaderImg     string `json:"header_img"`
	NickName      string `json:"nick_name"`
	BlogWant      string `json:"blog_want"comment:"blog签名"`
}

type BlogCategoryTaglistResponse struct {
	Tag      []blog.Tag
	Category []blog.Category
}

type BlogArticleResponse struct {
	blog.Article
	IsGood bool `json:"is_good"`
}

type BlogArticleToday struct {
	Total int64 `json:"total"`
}
type BlogArticleTodayNumber struct {
	Rate float64 `json:"rate"`
}
