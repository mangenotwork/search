package http_service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/search/api"
	"github.com/mangenotwork/search/core"
	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils"
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
	theme := c.Param("theme")   // *theme  主题
	word := c.Query("word")     // *word  搜索词
	sortType := c.Query("sort") // sort  排序类型 t: 时间，  o: 排序值, f: 词频   默认t
	// *out   输出结构  默认 id
	// id: 只输出docId,
	// list: 输出列表有 docId title author time_stamp OrderInt ,
	// full: 输出除 content 外的数据，并且含有关键词的 位置信息数据
	out := c.Query("out")
	pgStr := c.Query("pg")       // pg  页数   默认 1
	countStr := c.Query("count") // count  每页是数量  最大值不超过 100   默认100

	if len(theme) < 1 {
		APIOutPut(c, 1, 0, "", "缺少参数 theme ")
		return
	}

	if len(word) < 1 {
		APIOutPut(c, 1, 0, "", "缺少搜索词 word ")
		return
	}

	if sortType != "t" && sortType != "o" && sortType != "f" {
		sortType = "t"
	}

	if out != "id" && out != "list" && out != "full" {
		out = "id"
	}

	pg := utils.Any2Int(pgStr)
	if pg < 1 {
		pg = 1
	}

	count := utils.Any2Int(countStr)
	if count < 1 {
		count = 100
	}

	switch out {
	case "id":
		data, err := new(api.APISearch).SearchId(theme, word, sortType, count)
		if err != nil {
			APIOutPut(c, 1, 0, "", err.Error())
			return
		}
		APIOutPut(c, 0, 0, data, "ok")
	case "list":
		data, err := new(api.APISearch).SearchList(theme, word, sortType, count)
		if err != nil {
			APIOutPut(c, 1, 0, "", err.Error())
			return
		}
		APIOutPut(c, 0, 0, data, "ok")
	case "full":
		data, err := new(api.APISearch).SearchFull(theme, word, sortType, count)
		if err != nil {
			APIOutPut(c, 1, 0, "", err.Error())
			return
		}
		APIOutPut(c, 0, 0, data, "ok")
	}

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

func TermData(c *gin.Context) {
	theme := c.Param("theme") //主题
	word := c.Query("word")
	sortType := c.Query("sort") // t: 时间，  o: 排序值, f: 词频
	data := new(api.APISearch).GetTermData(theme, word, sortType, 100)
	APIOutPut(c, 0, 0, data, "ok")
}
