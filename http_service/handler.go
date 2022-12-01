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
	theme := c.Param("theme") //主题
	logger.Info("theme = ", theme)

	param := &entity.Doc{}
	err := c.BindJSON(param)
	if err != nil {
		APIOutPut(c, 1, 0, "", err.Error())
		return
	}
	new(api.APIDoc).Set(theme, param)
	APIOutPut(c, 0, 0, "ok", "ok")
}

func GetDoc(c *gin.Context) {
	theme := c.Param("theme") //主题
	docId := c.Query("doc_id")
	data, err := new(api.APIDoc).Get(theme, docId)
	if err != nil {
		APIOutPut(c, 1, 0, "", err.Error())
		return
	}
	APIOutPut(c, 0, 0, data, "ok")
}

func CreatedTheme(c *gin.Context) {
	param := &entity.Theme{}
	err := c.BindJSON(param)
	if err != nil {
		APIOutPut(c, 1, 0, "", err.Error())
		return
	}
	err = new(api.APITheme).Created(param)
	if err != nil {
		APIOutPut(c, 1, 0, "", err.Error())
		return
	}
	APIOutPut(c, 0, 0, "", "创建成功")
}

func GetThemeList(c *gin.Context) {
	data := new(api.APITheme).GetAll()
	APIOutPut(c, 0, 0, data, "ok")
}

func GetTheme(c *gin.Context) {
	name := c.Query("name")
	data, err := new(api.APITheme).Get(name)
	if err != nil {
		APIOutPut(c, 1, 0, "", err.Error())
		return
	}
	APIOutPut(c, 0, 0, data, "ok")
}

func DelTheme(c *gin.Context) {

}

func UpdateTheme(c *gin.Context) {

}
