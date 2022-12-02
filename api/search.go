package api

import (
	"github.com/mangenotwork/search/core"
	"github.com/mangenotwork/search/entity"
)

type APISearch struct {
}

func (api *APISearch) GetTermData(word, sortType string) []*entity.PL {
	return core.GetSearchFile(word, sortType)
}
