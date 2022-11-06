package blog

import "blog/global"

type LikeAndWatch struct {
	global.GM_MODEL
	Ip        string `json:"ip"`
	ArticleID uint   `json:"article_id"`
	//1 true 2 false
	Like uint8 `json:"like"gorm:"default:0"gorm:"default:0;comment:是否有帮助"`
}
