package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// SubmitMaterialAudit 提交素材审核列表
func SubmitMaterialAudit(db *gorm.DB, auditList types.MaterialAuditList) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, audit := range auditList {
			if err := submitMaterialAudit(db, audit); err != nil {
				return err
			}
		}
		return nil
	})
}

func submitMaterialAudit(db *gorm.DB, audit *types.MaterialAudit) error {
	if err := Upsert(db, audit); err != nil {
		return err
	}

	if err := db.Model(&types.Material{ID: audit.MaterialID}).
		Updates(map[string]interface{}{"AuditStatus": audit.AuditStatus, "RejectReason": audit.RejectReason}).Error; err != nil {
		return err
	}

	return nil
}

// ListMaterialAudit 列出所有素材审核, 包含素材、审核人信息
func ListMaterialAuditDetail(db *gorm.DB, option *StatementOption) ([]*types.MaterialAudit, error) {
	audits := make([]*types.MaterialAudit, 0)
	if err := getDBWithOption(db, option).Joins("Material").Joins("User").Find(&audits).Error; err != nil {
		return nil, err
	}

	return audits, nil
}
