package handler

import (
    "fmt"

    "github.com/tencentad/martech/api/types"
    "github.com/tencentad/martech/pkg/orm"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// BackendGetHandler 拉取页面实体列表
func BackendGetHandler(c *gin.Context) {
    dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
        return orm.GetAllBackend(db)
    })
}

// BackendPostHandler 新增页面实体信息
func BackendPostHandler(c *gin.Context) {
    item := &types.Backend{}
    dbEditHandler(c, item, func(db *gorm.DB) error {
        name := fmt.Sprint(item.Path, "_", item.Method)
        // old
        if item.ID != 0 {
            old := &types.Backend{ID: item.ID}
            if err := orm.TakeBackend(db, old); err != nil {
                return err
            }

            if item.ObjectID != old.ObjectID {
                return fmt.Errorf("can't update object_id")
            }

            // update object and backend
            return db.Transaction(func(tx *gorm.DB) error {
                item.Object = &types.Object{
                    ID:   item.ObjectID,
                    Name: name,
                }

                if err := orm.UpsertBackend(tx, item); err != nil {
                    return err
                }
                return orm.UpsertObject(tx, item.Object)
            })
        }

        // create new one
        return db.Transaction(func(tx *gorm.DB) error {
            o := &types.Object{
                Name: name,
                Type: types.ObjectTypeBackend,
                TenantID: SuperAdminDomainID,
            }
            item.Object = o

            return orm.UpsertBackend(tx, item)
        })
    })
}

// BackendDeleteHandler 删除某个页面实体
func BackendDeleteHandler(c *gin.Context) {
    dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
        item := &types.Backend{ID: id}
        if err := orm.TakeBackend(db, item); err != nil {
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
            return orm.DeleteBackendById(db, id)
        })
    })
}

