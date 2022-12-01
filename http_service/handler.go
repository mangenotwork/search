package http_service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/search/api"
	"github.com/mangenotwork/search/core"
	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils/logger"
	"net/http"
	"time"
)

// ResponseJson 统一接口输出
type ResponseJson struct {
	Code      int64       `json:"code"`
	Page      int64       `json:"page"`
	Msg       string      `json:"msg"`
	Body      interface{} `json:"body"`
	Take      int64       `json:"took"`
	TakeStr   string      `json:"took_str"`
	TimeStamp int64       `json:"timeStamp"`
}

// APIOutPut 统一接口输出方法
func APIOutPut(c *gin.Context, code int64, count int, data interface{}, msg string) {
	// TODO 通过 count 计算页
	tum, _ := c.Get("tum")
	t2 := time.Now().UnixNano()
	t := t2 - tum.(int64)
	resp := &ResponseJson{
		Code:      code,
		Msg:       msg,
		Body:      data,
		TimeStamp: time.Now().Unix(),
		Take:      t,
		TakeStr:   fmt.Sprintf("%fms", float64(t)/1e6),
	}
	c.IndentedJSON(http.StatusOK, resp)
	return
}

func Index(c *gin.Context) {

	core.Case()

	APIOutPut(c, 0, 0, "ok", "ok")
}

func Search(c *gin.Context) {

	rse := core.Search()

	APIOutPut(c, 0, 0, rse, "ok")
}

func GetTerm(c *gin.Context) {
	str := c.Query("str")
	rse := new(api.APIFenCi).TermExtract(str)
	APIOutPut(c, 0, 0, rse, "ok")
}

func SetDoc(c *gin.Context) {
	theme := c.Query("theme")       //主题
	isAuthor := c.Query("author")   // 是否对 author 创建索引 0:否，1:是
	isContent := c.Query("content") // 是否对 content 创建索引 0:否，1:是

	logger.Info("theme = ", theme)
	logger.Info("author = ", isAuthor)
	logger.Info("content = ", isContent)

	param := &entity.Doc{}
	err := c.BindJSON(param)
	if err != nil {
		APIOutPut(c, 1, 0, "", err.Error())
		return
	}
	new(api.APIDoc).Set(param)
	APIOutPut(c, 0, 0, "ok", "ok")
}

func GetDoc(c *gin.Context) {
	docId := c.Query("doc_id")
	data, err := new(api.APIDoc).Get(docId)
	if err != nil {
		APIOutPut(c, 1, 0, "", err.Error())
		return
	}
	APIOutPut(c, 0, 0, data, "ok")
}

func NewTheme(c *gin.Context) {
	APIOutPut(c, 0, 0, "", "ok")
}
