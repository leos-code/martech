package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RtaExpListHandler 获取Rta实验列表的接口
func RtaExpListHandler(c *gin.Context) {
	account := &types.RtaAccount{}
	dbGetHandler(c, account, func(db *gorm.DB) (interface{}, error) {
		err := orm.LoadAccountRtaExp(db, account)
		return account, err
	})
}
