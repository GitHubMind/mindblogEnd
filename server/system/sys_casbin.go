package server

import (
	"blog/global"
	"blog/model/system/request"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

type CasbinService struct {
}

var CasbinServiceApp = new(CasbinService)

func (casbinService CasbinService) UpdateCasbin(AuthorityID uint, casbinInfos []request.CasbinInfo) error {
	authorityId := strconv.Itoa(int(AuthorityID))
	//  先删除 再添加
	casbinService.ClearCasbin(0, authorityId)
	rules := [][]string{}
	for _, v := range casbinInfos {
		rules = append(rules, []string{authorityId, v.Path, v.Method})
	}
	//再去添加
	if success := casbinService.AddCasbin(rules); !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

//@function: ClearCasbin
//@description: 清除匹配的权限
//@param: v int, p ...string
//@return: bool
func (casbinService *CasbinService) AddCasbin(rules [][]string) bool {
	//单利创建
	e := casbinService.Casbin()
	//删除改 csbin policy
	success, _ := e.AddPolicies(rules)
	return success
}

//@function: AddCasbin
//@description: 添加匹配的权限
//@param: v int, p ...string
//@return: bool

func (casbinService *CasbinService) ClearCasbin(v int, p ...string) bool {
	//单利创建
	e := casbinService.Casbin()
	//删除改 csbin policy
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

//@function: GetPolicyPathByAuthorityId
//@description: 获取权限列表
//@param: authorityId string
//@return: pathMaps []request.CasbinInfo

func (casbinService *CasbinService) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []request.CasbinInfo) {
	//经典单利
	e := casbinService.Casbin()
	authorityId := strconv.Itoa(int(AuthorityID))
	//filedindex等于偏移量  例如v0 是存储id的一样
	list := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, request.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

//@function: Casbin
//@description: 持久化到数据库  引入自定义规则
//@return: *casbin.Enforcer

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func (casbinService *CasbinService) Casbin() *casbin.SyncedEnforcer {
	//单例
	once.Do(func() {
		//驱动器 用mysql 然后他自己会直接写定mysql
		a, _ := gormadapter.NewAdapterByDB(global.GM_DB)
		//perm原则写入
		text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		//match 方式很湿很粗暴的 三个对三个 具体难一点对可能要另外分析
		m, err := model.NewModelFromString(text)
		if err != nil {
			zap.L().Error("字符串加载模型失败!", zap.Error(err))
			return
		}
		syncedEnforcer, _ = casbin.NewSyncedEnforcer(m, a)
	})
	_ = syncedEnforcer.LoadPolicy()
	return syncedEnforcer
}

//@function: UpdateCasbinApi
//@description: API更新随动
//@param: oldPath string, newPath string, oldMethod string, newMethod string
//@return: error

func (casbinService *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	//
	err := global.GM_DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}
