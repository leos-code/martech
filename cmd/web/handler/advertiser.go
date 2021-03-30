package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tencentad/union-marketing-go-sdk/api/manager"
	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	authURIInputStataFormat = "{\"platform\":\"%v\",\"object_id\":__OBJECT_ID__}"
)

type stateData struct {
	Platform sdk.MarketingPlatformType `json:"platform"`
	ObjectID uint64                    `json:"object_id"`
}

// AdvertiserAuthorizeHandler 拉取授权列表
func AdvertiserAuthorizeHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		list := types.AdvertiserAuthorizations{}
		platforms := manager.GetPlatformList()
		for _, v := range platforms {
			impl, err := manager.GetImpl(v)
			if err != nil {
				return nil, err
			}

			state := fmt.Sprintf(authURIInputStataFormat, v)
			input := &sdk.GenerateAuthURIInput{State: state}
			output, err := impl.GenerateAuthURI(input)
			if err != nil {
				return nil, err
			}

			list = append(list, &types.AdvertiserAuthorization{
				Platform: v,
				Url:      output.AuthURI,
			})
		}

		return list, nil
	})
}

// AdvertiserCallbackHandler 授权回调接口
func AdvertiserCallbackHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)

		input := c.DefaultQuery("state", "")
		state := &stateData{}
		if err := json.Unmarshal([]byte(input), state); err != nil {
			return nil, err
		}

		if state.ObjectID == 0 {
			return nil, ErrEmptyID
		}

		impl, err := manager.GetImpl(state.Platform)
		if err != nil {
			return nil, err
		}

		e := GetEnforcer(c)
		object := &types.Object{ID: state.ObjectID}
		if err := orm.TakeObject(db, object); err != nil {
			return nil, err
		}
		if object.Type != types.ObjectTypeAdvertiser {
			return nil, ErrNoWritePermission
		}

		oo := casbinObjectEncode(object)
		if ok, _ := e.Enforce(cu, domain, oo, write); !ok {
			return nil, ErrNoWritePermission
		}
		output, err := impl.ProcessAuthCallback(&sdk.ProcessAuthCallbackInput{
			AuthCallback: c.Request,
		})
		if err != nil {
			return nil, err
		}

		var advertisers []*types.Advertiser
		for _, v := range output.AuthAccountList {
			v.AccessToken = ""
			v.AccessTokenExpireAt = time.Time{}
			v.RefreshToken = ""
			v.RefreshTokenExpireAt = time.Time{}

			u := types.AuthAccount(*v)
			advertiser := &types.Advertiser{
				Platform: v.Platform,
				Name:     v.ID,
				TenantID: currentT.ID,
				ObjectID: state.ObjectID,
				AuthAccount: &u,
			}

			advertisers = append(advertisers, advertiser)
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			for _, v := range advertisers {
				old := &types.Advertiser{
					Name:     v.Name,
					TenantID: currentT.ID,
				}
				// deleted or no exists Advertiser is no biyao to check write
				if err := orm.TakeAdvertiser(db, old); err == nil {
					o1 := &types.Object{ID: old.ObjectID}
					o2 := casbinObjectEncode(o1)
					if ok, _ := e.Enforce(cu, domain, o2, write); !ok {
						return ErrNoWritePermission
					}
				}

				if err := orm.Upsert(tx, v); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return nil, err
		}

		c.Redirect(http.StatusFound, "/page/proxy#advertiser")
		return nil, nil
	})
}

func AdvertiserGetHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		advertisers, err := orm.ListAdvertiserByTenantId(db, currentT.ID)
		if err != nil {
			return nil, err
		}

		linq.From(advertisers).Where(func(v interface{}) bool {
			advertiser := v.(*types.Advertiser)
			oo := casbinObjectEncode(&types.Object{ID: advertiser.ObjectID})
			ok, _ := e.Enforce(cu, domain, oo, read)
			return ok
		}).ToSlice(&advertisers)

		return orm.GetAllAdvertiser(db)
	})
}

func AdvertiserPatchHandler(c *gin.Context) {
	item := &types.Advertiser{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		// 1. 先根据id找到现在在数据库中的advertiser
		currentItem := &types.Advertiser{ID: item.ID}
		if err := orm.TakeAdvertiser(db, currentItem); err != nil {
			return err
		}
		// 2. 然后根据其objectId，判断是否有写现在object的权限
		oldObject := casbinObjectEncode(&types.Object{ID: currentItem.ObjectID})
		if ok, _ := e.Enforce(cu, domain, oldObject, write); !ok {
			return ErrNoWritePermission
		}
		// 3. 然后再查看是否有写现有的object的权限
		newObject := &types.Object{ID: item.ObjectID}
		if err := orm.TakeObject(db, newObject); err != nil {
			return err
		}
		if newObject.Type != types.ObjectTypeAdvertiser {
			return fmt.Errorf("incorrect object type")
		}
		no := casbinObjectEncode(newObject)
		if ok, _ := e.Enforce(cu, domain, no, write); !ok {
			return ErrNoWritePermission
		}
		// 4. 更新数据
		aq := &types.Advertiser{
			ID:       item.ID,
			ObjectID: item.ObjectID,
		}
		return orm.UpsertAdvertiser(db, aq)
	})
}

func AdvertiserDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		advertiser := &types.Advertiser{ID: id}
		if err := orm.TakeAdvertiser(db, advertiser); err != nil {
			return err
		}
		co := casbinObjectEncode(&types.Object{ID: advertiser.ObjectID})
		if ok, _ := e.Enforce(cu, domain, co, write); !ok {
			return ErrNoWritePermission
		}
		return orm.DeleteAdvertiser(db, advertiser.ID)
	})
}
