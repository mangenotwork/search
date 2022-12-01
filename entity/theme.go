package entity

// Theme 主题 可以理解为数据库的库
// 索引规则，默认对标题进行索引，强制
// 默认  0: false, 其他为true
type Theme struct {
	Name      string `json:"name"`
	IsAuthor  int    `json:"is_author"`  // 索引规则，是否对昵称进行索引，昵称不分词
	IsContent int    `json:"is_content"` // 索引规则，是否对内容文本进行索引，昵称不分词
}
