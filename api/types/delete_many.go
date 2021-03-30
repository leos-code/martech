package types

// DeleteMany 批量删除 POST Body
type DeleteMany struct {
	ID []uint64 `json:"id"` // 	待删除的ID列表
}
