package api

import (
	"bufio"
	"fmt"
	"github.com/mangenotwork/search/core"
	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils"
	"github.com/mangenotwork/search/utils/logger"
	"os"
)

type APIDoc struct {
}

// Set 写入文档
func (api *APIDoc) Set(theme string, doc *entity.Doc) {
	// 存储文档， 存在则更新
	docTheme := entity.DocPath + theme
	utils.Mkdir(docTheme)

	docFilePath := docTheme + "/" + doc.ID
	logger.Info("docFilePath = ", docFilePath)

	docData, err := utils.DataEncoder(doc)
	if err != nil {
		logger.Error("文档不能被压缩, err = ", err)
		return
	}
	file, err := os.OpenFile(docFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("文件打开失败", err)
		return
	}
	//及时关闭file句柄
	defer func() {
		_ = file.Close()
	}()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	_, err = write.Write(docData)
	if err != nil {
		logger.Error("文件写入失败, err = ", err)
	}
	// Flush将缓存的文件真正写入到文件中
	err = write.Flush()
	if err != nil {
		logger.Error("写入到文件中失败, err = ", err)
	}

	// TODO 创建倒排索引
	// 对title 强制倒排索引
	core.SetPostings(theme, doc.ID, doc.Title, doc.TimeStamp, doc.OrderInt)

}

func (api *APIDoc) Get(theme string, docId string) (*entity.Doc, error) {
	data := &entity.Doc{}
	docFilePath := entity.DocPath + theme + "/" + docId
	content, err := os.ReadFile(docFilePath)
	if err != nil {
		logger.Error("read file error:%v\n", err)
		return data, fmt.Errorf("no such doc.")
	}
	err = utils.DataDecoder(content, &data)
	return data, err
}
