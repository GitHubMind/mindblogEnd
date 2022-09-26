package lib

//有点害怕出现 cycle 的包导入
var (
	RegisterVerify         = Rules{"Username": {Ge("6")}, "NickName": {Ge("6")}, "Password": {Ge("6")}, "AuthorityId": {NotEmpty()}}
	LoginVerify            = Rules{"Username": {Ge("6")}, "Password": {Ge("6")}, "Captcha": {NotEmpty()}, "CaptchaId": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"CaptchaId": {NotEmpty()}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	MenuVerify             = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify         = Rules{"Title": {NotEmpty()}}
	IdVerify               = Rules{"ID": []string{NotEmpty()}}
	TitleVerify            = Rules{"Title": {NotEmpty()}}
	ArtileContentVerify    = Rules{"Content": {NotEmpty()}, "ID": []string{NotEmpty()}}
	ApiVerify              = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	ArticleVerify          = Rules{"title": {NotEmpty()}, "Desc": {NotEmpty()}, "Content": {NotEmpty()}, "CoverImageUrl": {NotEmpty()}}
	CreateArticleVerify    = Rules{"title": {NotEmpty()}, "Desc": {NotEmpty()}, "CoverImageUrl": {NotEmpty()}, "ArticleTags": {Ge("0")}, "Category": {Ge("0")}}
	TagVerify              = Rules{"name": {NotEmpty()}}
)
