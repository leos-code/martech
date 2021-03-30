package types

// RedirectLogin 单个跳转登录数据结构
type RedirectLogin struct {
	Type LoginType `json:"type"`
	Url  string    `json:"url"`
}

// RedirectLogins 跳转登录数据结构
type RedirectLogins []*RedirectLogin

// RegisterUserInfo 用户注册信息数据结构
type RegisterUserInfo struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email"        binding:"required"`
}
