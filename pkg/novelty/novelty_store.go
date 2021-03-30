package novelty

import "github.com/tencentad/martech/api/proto/novelty"

// UserNoveltyStore 获取用户新鲜度信息接口
type UserNoveltyStore interface{
	// Load
	Load(id string) (*novelty.Novelty, error)

	// Store 获取新鲜度数据
	Store(id string, data *novelty.Novelty) error
}
