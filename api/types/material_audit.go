package types

import (
	"time"

	"gorm.io/gorm"
)

// MaterialAudit 素材审核
type MaterialAudit struct {
	ID           uint64              `gorm:"column:id;primaryKey"       json:"id,omitempty"`
	CreatedAt    time.Time           `gorm:"column:created_at"          json:"created_at,omitempty"`
	UpdatedAt    time.Time           `gorm:"column:updated_at"          json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt      `gorm:"column:delete_at;index"     json:"-"`
	MaterialID   uint64              `gorm:"column:material_id"         json:"material_id,omitempty"`
	Material     *Material           `gorm:"foreignKey:MaterialID"      json:"material"`
	UserID       uint64              `gorm:"column:user_id"             json:"user_id,omitempty"`
	User         *User               `gorm:"foreignKey:UserID"          json:"user"`
	AuditStatus  MaterialAuditStatus `gorm:"column:audit_status"        json:"audit_status,omitempty"`
	RejectReason string              `gorm:"column:reject_reason"       json:"reject_reason,omitempty"`
}

type MaterialAuditStatus string

const (
	MaterialAuditUnaudited MaterialAuditStatus = "unaudited"
	MaterialAuditPass      MaterialAuditStatus = "pass"
	MaterialAuditReject    MaterialAuditStatus = "reject"
)

// MaterialAuditList  批量
type MaterialAuditList []*MaterialAudit
