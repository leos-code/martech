package orm

import (
	"github.com/tencentad/martech/api/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SchemaUpsert(db *gorm.DB, s *types.Schema) error {
	if s.ID == 0 {
		return db.Create(s).Error
	}else {
		return db.Updates(s).Error
	}
}

func SchemaGet(db *gorm.DB, c *gin.Context) (*types.Schema, error) {
	return nil, nil
}
