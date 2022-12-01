package entity

// Doc 文档
type Doc struct {
	ID          string `json:"id"`          // 文档的id, 唯一的，非唯一会导致文档被复写
	Title       string `json:"title"`       // 文档标题，默认对标题进行索引，不可取消
	Author      string `json:"author"`      // 文档作者，可选择对作者进行索引
	TimeStamp   int64  `json:"time_stamp"`  // 文档的日期(时间戳)， 对文档进行排序用
	OrderInt    int64  `json:"order_int"`   // 文档排序值，对文档进行排序
	Content     string `json:"content"`     // 文档内容
	Description string `json:"description"` // 文档描述
}
