package api

import (
	"fmt"
	"github.com/mangenotwork/search/core"
	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils/logger"
)

type APISearch struct {
}

func (api *APISearch) GetTermData(theme, word, sortType string, pg int) []*entity.PL {
	return core.GetSearchFile(theme, word, sortType, pg)
}

func (api *APISearch) SearchId(theme, word, sortType string, pg, count int) ([]string, error) {
	data := make([]string, 0)
	s, err := api.Search(theme, word, sortType, pg, count)
	if err != nil {
		return data, err
	}
	for _, v := range s {
		data = append(data, v.Key)
	}
	return data, nil
}

type OutList struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	TimeStamp int64  `json:"timestamp"`
	Author    string `json:"author"`
}

func (api *APISearch) SearchList(theme, word, sortType string, pg, count int) ([]*OutList, error) {
	data := make([]*OutList, 0)
	s, err := api.Search(theme, word, sortType, pg, count)
	if err != nil {
		return data, err
	}

	for _, v := range s {
		doc, _ := new(APIDoc).Get(theme, v.Key)
		data = append(data, &OutList{
			Id:        v.Key,
			Title:     doc.Title,
			TimeStamp: doc.TimeStamp,
			Author:    doc.Author,
		})
	}

	return data, nil
}

type FullList struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	TimeStamp   int64   `json:"timestamp"`
	Author      string  `json:"author"`
	OrderInt    int64   `json:"order_int"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	End         int     `json:"end"`
	Start       int     `json:"start"`
	Theme       string  `json:"theme"`
	SortType    string  `json:"sort_type"`
	SortValue   float64 `json:"sort_value"`
	Md5         string  `json:"md5"`
}

func (api *APISearch) SearchFull(theme, word, sortType string, pg, count int) ([]*FullList, error) {
	data := make([]*FullList, 0)
	s, err := api.Search(theme, word, sortType, pg, count)
	if err != nil {
		return data, err
	}

	for _, v := range s {
		// TODO 在缓存获取
		doc, _ := new(APIDoc).Get(theme, v.Key)
		data = append(data, &FullList{
			Id:          v.Key,
			Title:       doc.Title,
			TimeStamp:   doc.TimeStamp,
			Author:      doc.Author,
			OrderInt:    doc.OrderInt,
			Description: doc.Description,
			Content:     doc.Content,
			End:         v.End,
			Start:       v.Start,
			Theme:       v.TermText,
			SortType:    sortType,
			SortValue:   v.Value,
			Md5:         doc.SumMD5,
		})
	}

	return data, nil
}

func (api *APISearch) Search(theme, word, sortType string, pg, count int) ([]*entity.PLTerm, error) {
	data := make([]*entity.PLTerm, 0)

	// 判断是否存在 theme
	_, err := new(APITheme).ThemeCacheGet(theme)
	if err != nil {
		return data, fmt.Errorf("not theme.")
	}

	// TODO 判断 word 已经是一个 term , 如果是就不用 分词取 term

	// 拆分搜索词, 聚合搜索结果
	// 首词(第一个词)的权重最大， 依次取满数量  TODO 搜索词的权重
	termList := core.TermExtract(word)

	// TODO 如果 termList 取不到词，例如 输入 的, 就匹配一个相关的词
	if len(termList) < 1 {
		logger.Info("没有找到词的索引，匹配一个相关的")

	}

	hasMap := map[string]struct{}{}
	l := 0 // 计数
	for _, v := range termList {
		logger.Error("v.Text = ", v.Text)
		logger.Error("pg = ", pg)
		pl := core.GetSearchFile(theme, v.Text, sortType, pg)

		// 第一个词就满足了条数，直接返回
		// 100是默认值
		if len(pl) >= 100 {
			for _, d := range pl {
				data = append(data, &entity.PLTerm{v.Text, *d})
			}
			return data, nil
		}

		// 第一个词不够，有好多往里填好多
		for _, d := range pl {
			if l >= 100 {
				// 数量够 直接返回
				return data, nil
			}
			_, ok := hasMap[d.Key]
			if !ok {
				hasMap[d.Key] = struct{}{}
				data = append(data, &entity.PLTerm{v.Text, *d})
				l++
			}
		}
	}
	return data, nil
}
