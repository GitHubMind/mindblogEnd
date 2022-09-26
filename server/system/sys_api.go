package server

import (
	"blog/global"
	"blog/model/commond/request"
	"blog/model/system"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ApiService struct{}

var ApiServiceApp = new(ApiService)

//@function: GetAPIInfoList
//@description: 分页获取数据,
//@param: api model.SysApi, info request.PageInfo, order string, desc bool
//@return: list interface{}, total int64, err error
func (apiService *ApiService) GetAPIInfoList(api system.SysApi, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	//分页,获取表的所有数据
	db := global.GM_DB.Debug().Model(&system.SysApi{})
	var apiList []system.SysApi

	if api.Path != "" {
		db = db.Where("path LIKE ?", "%"+api.Path+"%")
	}

	if api.Description != "" {
		db = db.Where("description LIKE ?", "%"+api.Description+"%")
	}

	if api.Method != "" {
		db = db.Where("method = ?", api.Method)
	}

	if api.ApiGroup != "" {
		db = db.Where("api_group = ?", api.ApiGroup)
	}
	//获取总数
	err = db.Count(&total).Error

	if err != nil {
		return apiList, total, err
	} else {
		//下一步添加分页数据
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 有效制定了，不过这个写法我感觉可以优化一下
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["path"] = true
			orderMap["api_group"] = true
			orderMap["description"] = true
			orderMap["method"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return apiList, total, err
			}
			//有排序的写法
			err = db.Order(OrderStr).Find(&apiList).Error
		} else {
			//无排序写法
			err = db.Order("api_group").Find(&apiList).Error
		}
	}
	return apiList, total, err
}

func (apiService *ApiService) GetAllApis() (apis []system.SysApi, err error) {
	//很直接就全盘拖出来
	err = global.GM_DB.Find(&apis).Error
	return
}

func (apiService *ApiService) CreateApi(api system.SysApi) (err error) {
	if !errors.Is(global.GM_DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同api")
	}
	return global.GM_DB.Create(&api).Error
}

//@description: 根据id获取api
//@param: id float64
//@return: api model.SysApi, err error

func (apiService *ApiService) GetApiById(id int) (api system.SysApi, err error) {
	err = global.GM_DB.Where("id = ?", id).First(&api).Error
	return
}

//@function: UpdateApi
//@description: 根据id更新api
//@param: api model.SysApi
//@return: err error

func (apiService *ApiService) UpdateApi(api system.SysApi) (err error) {
	var oldA system.SysApi
	//拿一个对比
	err = global.GM_DB.Where("id = ?", api.ID).First(&oldA).Error
	if oldA.Path != api.Path || oldA.Method != api.Method {
		if !errors.Is(global.GM_DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在相同api路径")
		}
	}
	if err != nil {
		return err
	} else {
		err = CasbinServiceApp.UpdateCasbinApi(oldA.Path, api.Path, oldA.Method, api.Method)
		if err != nil {
			return err
		} else {
			err = global.GM_DB.Save(&api).Error
		}
	}
	return err
}

//@description: 删除选中API
//@param: apis []model.SysApi
//@return: err error

func (apiService *ApiService) DeleteApisByIds(ids request.IdsReq) (err error) {
	//不同步删除 casbin吗,但是源码缺正常删除了、
	var entity []system.SysApi
	if errors.Is(global.GM_DB.Debug().Where("id in (?) ", ids.Ids).Find(&entity).Error, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "无数据可以删除")
	}

	err = global.GM_DB.Delete(&entity).Error
	for _, value := range entity {
		CasbinServiceApp.ClearCasbin(1, value.Path, value.Method)
	}

	return err
}

func (apiService *ApiService) DeleteApi(api system.SysApi) (err error) {
	var entity system.SysApi
	global.GM_DB.Delete(&api)
	if !errors.Is(global.GM_DB.Where("id =? ", api.ID).First(&entity).Error, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "无数据可以删除")
	}
	if global.GM_DB.Delete(&entity).Error != nil {
		return err
	}
	CasbinServiceApp.ClearCasbin(1, entity.Path, entity.Method)
	return nil
}
