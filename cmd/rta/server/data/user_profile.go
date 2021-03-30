package data

// UserProfile 用户画像
type UserProfile struct {
	Feature map[string]interface{}
}

func newUserProfile() *UserProfile {
	return &UserProfile{
		Feature: make(map[string]interface{}),
	}
}