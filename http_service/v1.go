package http_service

func V1(id string) {
	v1 := Router.Group("", HttpMiddleware(id))

	v1.GET("/", Index)

	// ============================================================ 测试接口

	v1.GET("/search/:theme/term/data", TermData) // 获取倒排结果

	// ============================================================ 接口
	v1.GET("/term", GetTerm)                       // 分词
	v1.POST("/theme", CreatedTheme)                // Theme 创建主题
	v1.GET("/theme/list", GetThemeList)            // Theme 查看主题列表
	v1.GET("/theme", GetTheme)                     // Theme 查看主题
	v1.DELETE("/theme", DelTheme)                  // TODO  Theme 删除主题
	v1.PUT("/theme", UpdateTheme)                  // TODO  Theme 修改主题
	v1.POST("/doc/:theme", SetDoc)                 // 写文件
	v1.GET("/doc/:theme/:doc_id", GetDoc)          // 读文档
	v1.GET("/search/:theme", Search)               // 搜索
	v1.GET("/doc/:theme/:doc_id/term", GetDocTerm) // 查看文档的 term
}
