package types

import (
	"time"

	"gorm.io/gorm"
)

// User 用户信息
type User struct {
	ID              uint64             `gorm:"column:id;primaryKey"             json:"id,omitempty"`
	CreatedAt       time.Time          `gorm:"column:created_at"                json:"created_at,omitempty"`
	UpdatedAt       time.Time          `gorm:"column:updated_at"                json:"updated_at,omitempty"`
	DeletedAt       gorm.DeletedAt     `gorm:"column:delete_at;index"           json:"-"`
	PhoneNumber     string             `gorm:"column:phone_number;unique"       json:"phone_number,omitempty"`
	Email           string             `gorm:"column:email;unique"              json:"email,omitempty"`
	LoginUser       []*LoginUser       `gorm:"foreignKey:UserID"                json:"login_user,omitempty"`
	ExperimentGroup []*ExperimentGroup `gorm:"many2many:user_experiment_groups" json:"experiment_group,omitempty"`
	Role            []*Role            `gorm:"-"                                json:"role,omitempty"`
}

// UserInfo 用户信息数据结构
type UserInfo struct {
	User          *User      `json:"user"`
	LoginUser     *LoginUser `json:"login_user"`
	CurrentTenant *Tenant    `json:"current_tenant"`
	Tenant        []*Tenant  `json:"tenant"`
}

// UserSearch 用户搜索数据结构
type UserSearch struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	NickName    string `json:"nick_name"`
	Approximate bool   `json:"approximate"`
}

// IsEmpty 判断用户搜索数据是否非法
func (u *UserSearch) IsEmpty() bool {
	return u.PhoneNumber == "" && u.Email == "" && u.NickName == ""
}
