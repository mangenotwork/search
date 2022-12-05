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
func (api *APIDoc) Set(theme string, doc *entity.Doc) error {
	themeObj, err := new(APITheme).ThemeCacheGet(theme)
	if err != nil {
		logger.Error("not theme")
		return fmt.Errorf("not theme")
	}

	// 存储文档， 存在则更新
	docTheme := entity.DocPath + theme
	utils.Mkdir(docTheme)

	docFilePath := docTheme + "/" + doc.ID
	logger.Info("docFilePath = ", docFilePath)

	docData, err := utils.DataEncoder(doc)
	if err != nil {
		logger.Error("文档不能被压缩, err = ", err)
		return err
	}
	file, err := os.OpenFile(docFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("文件打开失败", err)
		return err
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
		return err
	}
	// Flush将缓存的文件真正写入到文件中
	err = write.Flush()
	if err != nil {
		logger.Error("写入到文件中失败, err = ", err)
		return err
	}

	// 创建倒排索引
	// 对title 强制倒排索引
	core.SetPostings(theme, doc.ID, doc.Title, doc.TimeStamp, doc.OrderInt)

	// 如果设置了对名称倒排
	if themeObj.IsAuthor != 0 {
		logger.Info("对名称倒排")
		core.SetPostingsAuthor(theme, doc.ID, doc.Author, doc.TimeStamp, doc.OrderInt)
	}

	logger.Info("themeObj.IsContent = ", themeObj.IsContent)
	// 如果设置了对文本倒排
	if themeObj.IsContent != 0 {
		logger.Info("对文本倒排")
		core.SetPostings(theme, doc.ID, doc.Content, doc.TimeStamp, doc.OrderInt)
	}
	return err
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
