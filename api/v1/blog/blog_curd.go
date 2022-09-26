package blog

import (
	"blog/global"
	"blog/lib"
	"blog/model/blog"
	"blog/model/blog/request"
	"blog/model/commond/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type blogCurd struct {
}

// @Tags Article
// @Summary 创建文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ArticleRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/CreateArticle [post]
func (g blogCurd) CreateArticle(c *gin.Context) {
	var r blog.Article
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.CreateArticleVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	r.Uid = lib.GetUserID(c)
	id, err := blogService.CreateArticle(&r)
	if err != nil {
		response.FailWithMessage("添加失败,不能重复添加", c)
		return
	}
	response.OkWithDetailed(map[string]uint{"id": id}, "添加成功", c)
}

// @Tags Article
// @Summary 删除文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ArticleRequsetID true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/CreateArticle [delete]
func (g blogCurd) DeleteArticle(c *gin.Context) {
	var r request.ArticleRequsetID
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	id, err := blogService.DeleteArticle(&r)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		global.GM_LOG.Error(err)
		return
	}
	response.OkWithDetailed(map[string]uint{"id": id}, "删除成功", c)
} // @Tags Article
// @Summary 删除文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ArticleRequsetIDMul true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/DeleteMulArticle [delete]
func (g blogCurd) DeleteMulArticle(c *gin.Context) {
	var r request.ArticleRequsetIDMul
	_ = c.ShouldBindJSON(&r)
	id, err := blogService.DeleteMulArticle(r)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		global.GM_LOG.Error(err)
		return
	}
	response.OkWithDetailed(map[string]uint{"id": id}, "删除成功", c)
}

// @Tags Article
// @Summary 修改文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body blog.Article true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/UpdateArticle [post]
func (g blogCurd) UpdateArticle(c *gin.Context) {
	var r blog.Article
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := blogService.UpdateArticle(&r)
	if err != nil {
		response.FailWithMessage("修改失败", c)
		global.GM_LOG.Error(err)
		return
	}
	response.OkWithDetailed(map[string]uint{}, "修改成功", c)
}

// @Tags Article
// @Summary 修改文章内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ArticleContentBackUp true "文章内容"
// @Success 200 {object} response.Response{msg=string} "1"
// @Router /blog/UpdateArticle [post]
func (g blogCurd) UpdateArticleContent(c *gin.Context) {
	var r request.ArticleContent
	var article blog.ArticleContentBackUp
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	article.Content = r.Content
	article.ArticleID = r.ID
	err := blogService.UpdateArticleContent(&article)
	if err != nil {
		response.FailWithMessage("修改失败", c)
		global.GM_LOG.Error(err)
		return
	}
	response.OkWithDetailed(map[string]uint{}, "修改成功", c)
}

// @Tags Article
// @Summary 修改线上展示的内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ArticleContentBackUp true "文章内容"
// @Success 200 {object} response.Response{msg=string} "1"
// @Router /blog/UpdateArticleContentOnLine [post]
func (g blogCurd) UpdateArticleContentOnLine(c *gin.Context) {
	var r request.ArticleContent
	var article blog.Article
	var articleContent blog.ArticleContentBackUp

	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	articleContent.Content = r.Content
	articleContent.ArticleID = r.ID
	article.ID = r.ID
	article.Content = r.Content
	err := blogService.UpdateArticle(&article)
	err = blogService.UpdateArticleContent(&articleContent)
	if err != nil {
		response.FailWithMessage("修改失败", c)
		global.GM_LOG.Error(err)
		return
	}
	response.OkWithDetailed(map[string]uint{}, "修改成功", c)
}

// @Tags Article
// @Summary 通过id获取文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ArticleRequsetID true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/FindArticle [post]
func (g blogCurd) FindArticle(c *gin.Context) {
	var r request.ArticleRequsetID
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	value, err := blogService.FindArticle(&r)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		global.GM_LOG.Error(err)
		return
	}
	response.OkWithDetailed(&value, "查询成功", c)
}

// @Tags Article
// @Summary 获得文章列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.TagRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/GetSearchArticleList [get]
func (g blogCurd) GetSearchArticleList(c *gin.Context) {
	//应该是一个分页操作了
	var r request.ArticleSearchRequset
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	r.ID = lib.GetUserID(c)
	value, total, err := blogService.GetSearchArticleList(&r)
	if err != nil {
		response.FailWithMessage("查询失败,不能重复添加", c)
		global.GM_LOG.Error("GetSearchArticleList查询失败", err)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     value,
		Total:    total,
		Page:     r.Page,
		PageSize: r.PageSize,
	}, "获取成功", c)
}

// @Tags Article
// @Summary 获得博客文章列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.TagRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/GetSearchArticleList [get]
func (g blogCurd) GetBlogSearchArticleList(c *gin.Context) {
	//应该是一个分页操作了
	var r request.ArticleSearchRequset
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	value, total, err := blogService.GetSearchArticleList(&r)
	if err != nil {
		response.FailWithMessage("查询失败,不能重复添加", c)
		global.GM_LOG.Error("GetSearchArticleList查询失败", err)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     value,
		Total:    total,
		Page:     r.Page,
		PageSize: r.PageSize,
	}, "获取成功", c)
}

// @Tags Article
// @Summary 创建tag
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.TagRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string,data=blog.Tag} "返回信息"
// @Router /blog/CreateTag [post]
func (g blogCurd) CreateTag(c *gin.Context) {
	var r request.TagRequset
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.TagVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//写入数据库
	//var tag blog.Tag

	tag := blog.Tag{Name: r.Name, Uid: lib.GetUserID(c)}
	err := blogService.CreateTag(&tag)
	if err != nil {
		response.FailWithMessage("添加失败,不能重复添加", c)
		return
	}
	response.OkWithMessage("添加成功", c)

	//返回
}

// @Tags Article
// @Summary 查询tag
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "tag列表"
// @Router /blog/GetTaglist [get]
func (g blogCurd) GetTag(c *gin.Context) {
	tags, err := blogService.GetTag()
	if err != nil {
		global.GM_LOG.Error("查失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(tags, c)

}

// @Tags Article
// @Summary 删除tag
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.TagIdRequset true "类型id"
// @Success 200 {object} response.Response{msg=string} "状态"
// @Router /blog/DeleteTag [delete]
func (g blogCurd) DeleteTag(c *gin.Context) {
	var r request.TagIdRequset
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.TagVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//写入数据库
	//var tag blog.Tag
	err := blogService.DeleteTag(r.ID)
	if err != nil {
		response.FailWithMessage("添加失败,不能重复添加", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

//
////创建Category
//
// @Tags Article
// @Summary 创建Category
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.TagRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string,data=blog.Tag} "返回信息"
// @Router /blog/CreateCategory  [post]
func (g blogCurd) CreateCategory(c *gin.Context) {
	var r request.TagRequset
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.TagVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//写入数据库
	//var tag blog.Tag
	tag := blog.Category{Name: r.Name, Uid: lib.GetUserID(c)}
	err := blogService.CreateCategory(&tag)
	if err != nil {
		response.FailWithMessage("添加失败,不能重复添加", c)
		return
	}
	response.OkWithMessage("添加成功", c)

	//返回
}

// @Tags Article
// @Summary 创建category
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.TagRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/GetCategory  [get]
func (g blogCurd) GetCategory(c *gin.Context) {

	tags, err := blogService.GetCategory()
	if err != nil {

		global.GM_LOG.Error("查失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(tags, c)

}

// @Tags Article
// @Summary 删除category
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.TagRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/DeleteCategory [delete]
func (g blogCurd) DeleteCategory(c *gin.Context) {

	var r request.TagIdRequset
	_ = c.ShouldBindJSON(&r)
	//要先添加东西
	if err := lib.Verify(r, lib.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//n
	err := blogService.DeleteCategory(r.ID)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// @Tags blog
// @Summary 根据blog操作
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.BlogNikeNameRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/GetBlogInfoByName [get]
func (g blogCurd) GetBlogInfoByName(c *gin.Context) {
	var r request.BlogNikeNameRequset
	r.NickName = c.Query("nickname")
	//要先添加东西
	if err := lib.Verify(r, lib.TitleVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//n
	info, err := blogService.GetBlogInfoByName(&r)
	if err != nil {

		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(info, "查询成功", c)
}

// @Tags blog
// @Summary 根据blog操作
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.BlogNikeNameRequset true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /blog/GetBlogInfoByName [get]
func (g blogCurd) GetBlogInfoById(c *gin.Context) {
	var r blog.Article
	value, _ := strconv.ParseUint(c.Query("id"), 10, 32)
	r.ID = uint(value)
	//要先添加东西
	if err := lib.Verify(r, lib.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := blogService.GetBlogInfoById(&r)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(info, "查询成功", c)
}

// @Tags Article
// @Summary 查询tag
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "tag列表"
// @Router /blog/GetTaglist [get]
func (g blogCurd) GetBlogCategoryTaglistById(c *gin.Context) {
	//人的id
	value, _ := strconv.ParseUint(c.Query("id"), 10, 32)
	//info, err := blogService.GetBlogCategoryTaglistById(&r)
	if err != nil {
		global.GM_LOG.Error("查失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(tags, c)

}

// @Tags A
