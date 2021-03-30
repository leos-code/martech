package orm

type Page struct {
	Total int64       `json:"total"` // 总共的page数
	List  interface{} `json:"list"`  // 数据列表
}
