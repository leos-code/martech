package types

import (
	"gorm.io/gorm"
	"time"
)

// LoginUser 用户的登录信息
type LoginUser struct {
	OpenID    string         `gorm:"column:open_id;primaryKey"    json:"open_id,omitempty"`
	LoginType LoginType      `gorm:"column:login_type;primaryKey" json:"login_type,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"            json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"            json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"       json:"-"`
	UserID    uint64         `gorm:"column:user_id"               json:"user_id,omitempty"`
	Avatar    string         `gorm:"column:avatar"                json:"avatar,omitempty"`
	NickName  string         `gorm:"column:nick_name"             json:"nick_name,omitempty"`
}

// LoginType 登录类型
type LoginType string

const (
	LoginTypeRio    LoginType = "rio"
	LoginTypeQQ     LoginType = "qq"
	LoginTypeWechat LoginType = "wechat"
)
