package api

import "github.com/mangenotwork/search/core"

type APIFenCi struct {
}

// TermExtract 索引词提取
func (api *APIFenCi) TermExtract(str string) []string {
	return core.TermExtract(str)
}
