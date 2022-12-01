package core

import (
	"github.com/mangenotwork/search/utils/logger"
	"log"
	"os"
	"strings"
)

func Search() []string {
	work := "推特上发了一段视频"

	termList := TermExtract(work)

	termList = append(termList, work)

	log.Println("term list = ", termList)

	rse := NewSet()

	// 分开查
	for _, v := range termList {
		for _, a := range Find(v) {
			rse.Add(a)
		}
	}
	rseList := rse.List()
	logger.Info("rse = ", rseList)
	return rseList
}

func Find(term string) []string {
	invertedFile := "./index/" + term + "/1.t"
	content, err := os.ReadFile(invertedFile)
	if err != nil {
		logger.Error("read file error:%v\n", err)
		return []string{}
	}
	return strings.Split(string(content), ";")
}

// 可以用于去重
type Set map[string]struct{}

func NewSet() Set {
	return Set{}
}

func (s Set) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Delete(key string) {
	delete(s, key)
}

func (s Set) List() []string {
	l := make([]string, 0)
	for k, _ := range s {
		l = append(l, k)
	}
	return l
}
