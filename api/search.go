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

	// 拆分搜索词, 聚合搜索结果  TODO 多线程实现
	// 首词(第一个词)的权重最大， 取满数量
	termList := core.TermExtract(word)

	hasMap := map[string]struct{}{}
	l := 0 // TODO 翻页

	// TODO 需要并发
	for _, v := range termList {
		logger.Error("v.Text = ", v.Text)
		logger.Error("pg = ", pg)
		pl := core.GetSearchFile(theme, v.Text, sortType, pg)
		for i, d := range pl {
			if i > count {
				break
			}
			_, ok := hasMap[d.Key]
			if !ok {
				hasMap[d.Key] = struct{}{}
				data = append(data, &entity.PLTerm{v.Text, *d})
				l++
			}
		}
	}

	// TODO 更具 sortType 获取 doc 的数据

	return data, nil
}
