package core

import (
	"fmt"
	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils"
	"github.com/mangenotwork/search/utils/logger"
	"os"
	"strings"
)

func GetSearchFile(theme, term, sortTypeType string, pg int) []*entity.PL {
	filePath := entity.IndexPath + theme + "/" + term
	data := make([]*entity.PL, 0)
	//  t: 时间，  o: 排序值, f: 词频
	switch sortTypeType {
	case "t":
		filePath += fmt.Sprintf("/%d.plt", pg)
	case "o":
		filePath += fmt.Sprintf("/%d.plo", pg)
	case "f":
		filePath += fmt.Sprintf("/%d.plf", pg)
	}

	logger.Error("filePath = ", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("read file error:%v\n", err)
		return data
	}

	err = utils.DataDecoder(content, &data)
	if err != nil {
		logger.Error("解压数据失败 :%v\n", err)
		return data
	}

	logger.Error("data = ", data, len(data))

	return data
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
