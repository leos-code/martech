package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func DeveloperAppGetHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		return orm.GetAllDeveloperApp(db)
	})
}

func DeveloperAppPostHandler(c *gin.Context) {
	item := &types.DeveloperApp{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		currentT, err := currentTenantSessionGet(c)
		if err!= nil {
			return err
		}
		item.TenantID = currentT.ID
		hash, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		item.OAuthAccessToken = string(hash)
		return orm.UpsertDeveloperApp(db, item)
	})
}

func DeveloperAppDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
		return orm.DeleteDeveloperApp(db, id)
	})
}
