package blog

import (
	"blog/global"
	"blog/model/blog"
	"blog/model/blog/request"
	"blog/model/blog/response"
	"blog/model/system"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type BlogService struct {
}

func (receiver BlogService) CreateArticle(art *blog.Article) (id uint, err error) {

	err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
		art.ArticleContentBackUp.Content = "<h1> 内容</h1>"
		db := tx.Create(&art)
		id = art.ID
		return db.Error
	})
	return
}

//先放入update
func (receiver BlogService) UpdateArticle(art *blog.Article) (err error) {

	err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Debug().Updates(&art).Error
		return err
	})
	return
}

func (receiver BlogService) UpdateArticleContent(art *blog.ArticleContentBackUp) (err error) {

	err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
		//db := global.GM_DB.Create(&art)
		//id = art.ID
		//return db.Error
		return tx.Debug().Model(&art).Where("article_id = ? ", &art.ArticleID).UpdateColumns(&art).Error
	})
	return
}

func (receiver BlogService) FindArticle(art *request.ArticleRequsetID) (value blog.Article, err error) {

	var articles blog.Article
	articles.ID = art.ID
	err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
		//db := global.GM_DB.Create(&art)
		//id = art.ID
		//return db.Error
		err = tx.Model(blog.Article{}).Preload("Category").Preload("LikeAndWatchs").Preload("Tag").Joins("ArticleContentBackUp").First(&articles).Error
		value = articles
		return err
	})
	return
}
func (receiver BlogService) DeleteArticle(req *request.ArticleRequsetID) (id uint, err error) {

	var art blog.Article
	art.ID = req.ID
	id = req.ID
	err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
		//db:=global.GM_DB.Association("Tag").Delete(&art).Error(
		err := global.GM_DB.Model(&art).Association("Tag").Clear()
		if err != nil {
			return err
		}
		err = global.GM_DB.Model(&art).Association("Category").Clear()
		if err != nil {
			return err
		}
		err = global.GM_DB.Delete(&art).Error
		if err != nil {
			return err
		}
		//err := global.GM_DB.Table("blog_article_tag").Delete("")
		return nil
	})
	return
}
func (receiver BlogService) DeleteMulArticle(req request.ArticleRequsetIDMul) (id uint, err error) {

	var arts []blog.Article
	for _, value := range req.ID {
		var art blog.Article
		art.ID = value.ID
		arts = append(arts, art)
	}
	err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
		//db:=global.GM_DB.Association("Tag").Delete(&art).Error(
		err := tx.Model(&arts).Association("Tag").Clear()
		if err != nil {
			return err
		}
		err = tx.Model(&arts).Association("Category").Clear()
		if err != nil {
			return err
		}
		err = tx.Delete(&arts).Error
		if err != nil {
			return err
		}
		//err := global.GM_DB.Table("blog_article_tag").Delete("")
		return nil
	})
	return
}
func (receiver BlogService) GetSearchArticleList(art *request.ArticleSearchRequset) (articles []blog.Article, total int64, err error) {
	var article blog.Article
	limit := art.PageSize
	offset := art.PageSize * (art.Page - 1)
	db := global.GM_DB.Model(&article)

	if art.Title != "" {
		db = db.Where("title LIKE ?", art.Title+"%")
	}
	//
	if art.State == 2 || art.State == 1 {
		db = db.Where("state = ?", art.State)
	}

	if len(art.ArticleTags) != 0 {
		var array []blog.Tag
		for _, value := range art.ArticleTags {
			var a blog.Tag
			a.ID = value
			array = append(array, a)
		}
		var artId []int
		global.GM_DB.Select("id").Model(&array).Preload("Tag").Association("Article").Find(&artId)
		if len(artId) > 0 {
			db = db.Where(" id in (?) ", artId)
		}
	}
	if len(art.Category) != 0 {
		var array []blog.Category
		for _, value := range art.Category {
			var a blog.Category
			a.ID = value
			array = append(array, a)
		}
		var artId []int
		global.GM_DB.Select("id").Model(&array).Preload("Category").Association("Article").Find(&artId)
		if len(artId) > 0 {
			db = db.Where(" id in (?) ", artId)
		}
	}
	db = db.Where("uid = ? ", art.ID).Order("created_at desc")
	db.Count(&total)
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("LikeAndWatchs").Preload("Category").Preload("Tag").Find(&articles).Error
	if err != nil {
		return
	}
	return
}
func (receiver BlogService) CreateTag(tag *blog.Tag) (err error) {

	err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
		return global.GM_DB.Create(&tag).Error
	})
	return
}

func (receiver BlogService) CreateCategory(tag *blog.Category) (err error) {

	err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
		return global.GM_DB.Create(&tag).Error
	})
	return
}
func (receiver BlogService) GetCategory() (tags []blog.Category, err error) {
	db := global.GM_DB.Model(&tags).Find(&tags)
	err = db.Error
	return
}

func (receiver BlogService) DeleteCategory(id uint) (err error) {

	var tag blog.Category
	tag.ID = id
	err = global.GM_DB.Where(" id = ? ", id).Unscoped().Delete(&tag).Error
	return
}

func (receiver BlogService) GetTag() (tags []blog.Tag, err error) {

	db := global.GM_DB.Model(&tags).Find(&tags)
	err = db.Error
	return
}
func (receiver BlogService) DeleteTag(id uint) (err error) {

	var tag blog.Tag
	tag.ID = id
	err = global.GM_DB.Where(" id = ? ", id).Unscoped().Delete(&tag).Error
	return
}
func (receiver BlogService) GetBlogInfoByName(nikeName *request.BlogNikeNameRequset) (info response.BlogInfoResponse, err error) {
	var user system.SysUser
	user.NickName = nikeName.NickName
	// 因为这个字段他是unqiue的 不担心重复
	err = global.GM_DB.Model(&user).First(&user).Error
	if err != nil {
		return
	}
	global.GM_DB.Model(blog.Article{}).Where(" uid = ? ", user.ID).Count(&info.ArticleCount)
	global.GM_DB.Model(blog.Tag{}).Where(" uid = ? ", user.ID).Count(&info.TagCount)
	global.GM_DB.Model(blog.Category{}).Where(" uid = ? ", user.ID).Count(&info.CategoryCount)
	info.GitHubAddress = user.GitHubAddress
	info.HeaderImg = user.HeaderImg
	info.Uid = user.ID
	info.BlogWant = user.BlogWant

	return
}

func (receiver BlogService) GetBlogInfoById(art *blog.Article) (info response.BlogInfoResponse, err error) {
	var user system.SysUser
	err = global.GM_DB.Model(&art).First(&art).Error
	if err != nil {
		return
	}
	user.ID = art.Uid
	// 因为这个字段他是unqiue的 不担心重复
	err = global.GM_DB.Model(&user).First(&user).Error
	if err != nil {
		return
	}
	global.GM_DB.Model(blog.Article{}).Where(" uid = ? ", user.ID).Count(&info.ArticleCount)
	global.GM_DB.Model(blog.Tag{}).Where(" uid = ? ", user.ID).Count(&info.TagCount)
	global.GM_DB.Model(blog.Category{}).Where(" uid = ? ", user.ID).Count(&info.CategoryCount)
	info.GitHubAddress = user.GitHubAddress
	info.HeaderImg = user.HeaderImg
	info.Uid = user.ID
	info.NickName = user.NickName
	info.BlogWant = user.BlogWant
	return
}

func (receiver BlogService) GetBlogCategoryTaglistById(id uint) (info response.BlogCategoryTaglistResponse, err error) {
	//global.GM_DB.Model(blog.Article{}).Where(" uid = ? ", id).Find(&info.ArticleCount)
	global.GM_DB.Model(blog.Tag{}).Where(" uid = ? ", id).Find(&info.Tag)
	global.GM_DB.Model(blog.Category{}).Where(" uid = ? ", id).Find(&info.Category)
	return
}
func (receiver BlogService) ClickBlogLike(value *blog.LikeAndWatch) (err error) {
	//redis  articleid-ip-like
	key := string(value.ArticleID) + value.Ip + "like"
	//去重
	val, err := global.GM_REDIS.Incr(context.Background(), key).Result()
	if val > 1 {
		err = errors.New("重复点赞")
		return
	}

	value.Like = 1
	global.GM_DB.Debug().Model(&value).Where("ip = ?", value.Ip).Where("article_id = ? ", value.ArticleID).Updates(&value)

	//redis blog_article_like=1
	return
}
func (receiver BlogService) CanclBlogLike(value *blog.LikeAndWatch) (err error) {
	//redis  articleid-ip-like
	key := string(value.ArticleID) + value.Ip + "like"
	//去重
	_, err = global.GM_REDIS.Del(context.Background(), key).Result()
	if err != nil {
		return
	} else {
		go func() {
			value.Like = 2
			global.GM_DB.Debug().Model(&value).Where("ip = ?", value.Ip).Where("article_id = ? ", value.ArticleID).Updates(&value)
		}()
	}
	return
}
func (receiver BlogService) ClickBlog(value *blog.LikeAndWatch) (err error) {
	//redis  articleid-ip-like
	var total int64
	err = global.GM_DB.Model(&value).Where("ip = ?", value.Ip).Where("article_id = ? ", value.ArticleID).Find(&value).Count(&total).Error
	if err != nil {
		return
	}
	if total == 0 {
		global.GM_DB.Model(&value).Save(&value)
	}
	return
}

func (receiver BlogService) GetRateNumber(ID uint) (total int64, err error) {
	//redis  articleid-ip-like
	var value *blog.LikeAndWatch
	timeSelect := time.Now().Format("2006-01-02")

	t2, _ := time.ParseInLocation("2006-01-02", timeSelect, time.Local)
	var art []blog.Article
	// 因为这个字段他是unqiue的 不担心重复
	//提取 id
	err = global.GM_DB.Model(&art).Where("uid = ? ", ID).Find(&art).Error
	var article_id_arr []uint
	for _, artitem := range art {
		article_id_arr = append(article_id_arr, artitem.ID)
	}
	err = global.GM_DB.Model(&value).Where("created_at BETWEEN ? AND ? ", t2, t2.AddDate(0, 0, 1)).Where("article_id in (?) ", article_id_arr).Find(&value).Count(&total).Error
	if err != nil {
		return
	}
	if total == 0 {
		global.GM_DB.Model(&value).Save(&value)
	}
	return
}
func (receiver BlogService) GetRateLikeNumber(ID uint) (total float64, err error) {
	//redis  articleid-ip-like
	var value []blog.LikeAndWatch
	var art []blog.Article
	// 因为这个字段他是unqiue的 不担心重复
	//提取 id
	err = global.GM_DB.Model(&art).Where("uid = ? ", ID).Find(&art).Error
	var article_id_arr []uint
	for _, artitem := range art {
		article_id_arr = append(article_id_arr, artitem.ID)
	}
	var totalInt int64
	err = global.GM_DB.Model(blog.LikeAndWatch{}).Where("article_id in (?) ", article_id_arr).Find(&value).Count(&totalInt).Error
	like := 0

	if err != nil {
		return
	}
	for _, watch := range value {
		if watch.Like == 1 {
			like++
		}
		log.Println(watch)
	}
	total = float64(like) / float64(totalInt)

	return
}
