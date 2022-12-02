package api

import (
	"github.com/mangenotwork/search/core"
	"github.com/mangenotwork/search/entity"
)

type APIFenCi struct {
}

// TermExtract 索引词提取
func (api *APIFenCi) TermExtract(str string) []*entity.Term {
	return core.TermExtract(str)
}
