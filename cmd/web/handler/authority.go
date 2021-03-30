package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/web/ginx"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/ahmetb/go-linq/v3"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const (
	SessionID        = "session_id"
	CasbinDefaultKey = "casbin"
	CasbinModel      = "casbin_model"
	OahthAccessToken = "oauth_access_token"
	Timestamp        = "timestamp"
)

var (
	localCache = cache.New(2*time.Minute, 5*time.Minute)
)

// AuthMiddleware 登录验证和权限中间件
func AuthLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := userSessionGet(c); err != nil {
			_ = c.Error(err)
			c.Abort()
			c.Redirect(http.StatusFound, "/page/proxy#login")
			c.Abort()
		}
	}
}

// AuthAPIMiddleware 认证的中间件
func AuthAPIMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := authCheck(c, localCache); err != nil {
			ginx.ResponseWithError(c, http.StatusUnauthorized, err)
		}
	}
}

// DeveloperAppMiddleware developerApp认证的中间件
func DeveloperAppMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := tokenCheck(c); err != nil {
			ginx.ResponseWithError(c, http.StatusUnauthorized, err)
		}
	}
}

// EnforcerMiddleware casbin初始化权限中间件
func EnforcerMiddleware(e casbin.IEnforcer, model string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(CasbinDefaultKey, e)
		c.Set(CasbinModel, model)
	}
}

// GetEnforcer 获取enforcer
func GetEnforcer(c *gin.Context) casbin.IEnforcer {
	return c.MustGet(CasbinDefaultKey).(casbin.IEnforcer)
}

// CasbinModel 获取casbin model
func GetCasbinModel(c *gin.Context) string {
	return c.MustGet(CasbinModel).(string)
}

// GetObjects 获取当前用户当前域，指定类型且有权限的Object
func GetObjects(c *gin.Context, action string, objectType ...types.ObjectType) ([]*types.Object, error) {
	currentT, _ := currentTenantSessionGet(c)

	db := orm.GetDB()
	objects, err := orm.GetObjectInDomain(db, currentT.ID, objectType...)
	if err == nil {
		return nil, err
	}

	currentU, _ := userSessionGet(c)
	cu := casbinUserEncode(currentU)
	domain := casbinTenantEncode(currentT)
	e := GetEnforcer(c)

	// filter objects read permission
	linq.From(objects).Where(func(v interface{}) bool {
		ok, _ := e.Enforce(cu, domain, v.(*types.Object).Object, action)
		return ok
	}).ToSlice(&objects)
	return objects, nil
}

func GetObjectFilterOption(c *gin.Context, action string, objectType ...types.ObjectType) (*orm.FilterOption, error) {
	objects, err := GetObjects(c, action, objectType...)
	if err != nil {
		return nil, err
	}

	var oid []uint64
	for _, v := range objects {
		oid = append(oid, v.ID)
	}

	filter := &orm.FilterOption{
		Field:     "object_id",
		Operation: orm.FilterOperationIn,
		Value:     oid,
	}
	return filter, nil
}

func authCheck(c *gin.Context, cache *cache.Cache) error {
	currentU, _ := userSessionGet(c)
	currentT, err := currentTenantSessionGet(c)
	if err != nil {
		return err
	}

	cu := casbinUserEncode(currentU)
	domain := casbinTenantEncode(currentT)

	backend := &types.Backend{
		Path:   c.FullPath(),
		Method: c.Request.Method,
	}

	o := getBackendObject(cache, backend)
	e := GetEnforcer(c)
	if ok, _ := e.Enforce(cu, domain, o, read); !ok {
		return ErrNoAPIPermission
	}

	return nil
}

func getBackendObject(cache *cache.Cache, backend *types.Backend) string {
	key := backend.Path + "##" + backend.Method
	if _, ok := cache.Get(key); !ok {
		db := orm.GetDB()
		if err := orm.TakeBackend(db, backend); err != nil {
			cache.Set(key, "", 20*time.Second)
		} else {
			o := &types.Object{ID: backend.ObjectID}
			v := casbinObjectEncode(o)
			cache.SetDefault(key, v)
		}
	}

	v, _ := cache.Get(key)
	return v.(string)
}

func tokenCheck(c *gin.Context) error {
	token := c.Query(OahthAccessToken)
	timestamp := c.Query(Timestamp)

	nsec, err := strconv.Atoi(timestamp)
	if err != nil {
		return err
	}
	requestTime := time.Unix(0, int64(nsec))
	if err := timeCheck(requestTime); err != nil {
		return err
	}

	db := orm.GetDB()
	if err := orm.TakeDeveloperApp(db, &types.DeveloperApp{OAuthAccessToken: token}); err != nil {
		return ErrNoAPIPermission
	}
	return nil
}

func timeCheck(requestTime time.Time) error {
	currentTime := time.Now()
	if duration := currentTime.Sub(requestTime); duration.Minutes() > 10 {
		return fmt.Errorf("expired requet")
	}
	return nil
}
