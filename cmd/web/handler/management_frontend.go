package handler

import (
    "fmt"

    "github.com/tencentad/martech/api/types"
    "github.com/tencentad/martech/pkg/orm"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// FrontendGetHandler 拉取页面实体列表
func FrontendGetHandler(c *gin.Context) {
    dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
        return orm.GetAllFrontend(db)
    })
}

// FrontendPostHandler 新增页面实体信息
func FrontendPostHandler(c *gin.Context) {
    item := &types.Frontend{}
    dbEditHandler(c, item, func(db *gorm.DB) error {
        name := fmt.Sprint(item.Key, "_", item.Type)
        // old
        if item.ID != 0 {
            old := &types.Frontend{ID: item.ID}
            if err := orm.TakeFrontend(db, old); err != nil {
                return err
            }

            if item.ObjectID != old.ObjectID {
                return fmt.Errorf("can't update object_id")
            }

            // update object and frontend
            return db.Transaction(func(tx *gorm.DB) error {
                item.Object = &types.Object{
                    ID:   item.ObjectID,
                    Name: name,
                }

                if err := orm.UpsertFrontend(tx, item); err != nil {
                    return err
                }
                return orm.UpsertObject(tx, item.Object)
            })
        }

        // create new one
        return db.Transaction(func(tx *gorm.DB) error {
            o := &types.Object{
                Name: name,
                Type: types.ObjectTypeFrontend,
                TenantID: SuperAdminDomainID,
            }
            item.Object = o

            return orm.UpsertFrontend(db, item)
        })
    })
}

// FrontendDeleteHandler 删除某个页面实体
func FrontendDeleteHandler(c *gin.Context) {
    dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
        item := &types.Frontend{ID: id}
        if err := orm.TakeFrontend(db, item); err != nil {
            return err
        }

        domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
        object := casbinObjectEncode(&types.Object{ID: item.ObjectID})
        e := GetEnforcer(c)

        // delete object as item in current tenant
        os, _ := e.GetModel()["g"]["g2"].RM.GetRoles(object, domain)
        for _, v := range os {
            if _, err := e.RemoveNamedGroupingPolicy("g2", object, v, domain); err != nil {
                return err
            }
        }

        // delete object as group in current tenant
        if _, err := e.RemoveFilteredNamedGroupingPolicy("g2", 1, object, domain); err != nil {
            return err
        }

        return db.Transaction(func(tx *gorm.DB) error {
            if err := orm.DeleteObjectById(db, item.ObjectID); err != nil {
                return err
            }
            return orm.DeleteFrontendById(db, id)
        })
    })
}
