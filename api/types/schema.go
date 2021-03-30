package types

import (
	"database/sql/driver"
	"time"

	"github.com/tencentad/martech/pkg/schema"
	"gorm.io/gorm"
)

type Schema struct {
	ID        uint64         `gorm:"column:id;primaryKey"                   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"                      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"                      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"                 json:"-"`
	Object    string         `gorm:"column:object"                          json:"object,omitempty"`
	TenantID  uint64         `gorm:"column:tenant_id;index:idx_role,unique" json:"tenant_id,omitempty"`
	ParentID  uint64         `gorm:"-"                                      json:"parent_id"`
	Data      *SchemaData    `gorm:"column:data"                            json:"data"`
}

type SchemaData schema.Schema

// Scan SchemaData value into Jsonb, implements sql.Scanner interface
func (s *SchemaData) Scan(value interface{}) error {
	return scan(value, s)
}

// Value return json value, implement driver.Valuer interface
func (s SchemaData) Value() (driver.Value, error) {
	return value(s)
}
