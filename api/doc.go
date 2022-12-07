package api

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/mangenotwork/search/core"
	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils"
	"github.com/mangenotwork/search/utils/logger"
	"os"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type APIDoc struct {
}

func (api *APIDoc) SumMD5(doc *entity.Doc) (string, error) {
	b, err := json.Marshal(doc)
	if err != nil {
		logger.Error(err)
		return "", fmt.Errorf("非法数据")
	}
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Set 写入文档
func (api *APIDoc) Set(theme string, doc *entity.Doc) error {
	themeObj, err := new(APITheme).ThemeCacheGet(theme)
	if err != nil {
		logger.Error("not theme")
		return fmt.Errorf("not theme")
	}
	// 获取文档的一致性哈希
	hashId, err := core.NewMetaDataObj.Get(doc.ID)
	if err != nil {
		logger.Error("获取 theme Consistent hash id 失败 ", err)
	}
	logger.Info("hashId = ", hashId)

	// TODO 分配存储物理地址, 集群
	routersPath := entity.DocPath + theme + "/" + hashId
	backUpPath := entity.DocPath + theme + "/" + hashId + "_bk"
	utils.Mkdir(routersPath)
	utils.Mkdir(backUpPath)
	themeObj.MetaData[hashId] = &entity.MetaData{
		HashId:  hashId,
		Routers: routersPath,
		BackUp:  backUpPath,
	}

	docFilePath := routersPath + "/" + doc.ID
	logger.Info("docFilePath = ", docFilePath)
	docBackUpFilePath := backUpPath + "/" + doc.ID

	// 计算文档md5
	sumMD5, err := api.SumMD5(doc)
	if err != nil {
		return err
	}
	doc.SumMD5 = sumMD5
	docData, err := utils.DataEncoder(doc)
	if err != nil {
		logger.Error("文档不能被压缩, err = ", err)
		return err
	}

	// 查看是否存在，内容是否有变化
	if utils.FileExists(docFilePath) {
		docObj := &entity.Doc{}
		content, err := os.ReadFile(docFilePath)
		if err != nil {
			logger.Error(err)
		}
		err = utils.DataDecoder(content, &docObj)
		if docObj.SumMD5 == sumMD5 {
			// 数据一模一样什么都不干
			return fmt.Errorf("数据存在")
		}
	} else {
		// 新增数据
		themeObj.DocNumber++
	}

	// 写入文档
	err = api.doc2File(docFilePath, docData)
	if err != nil {
		return err
	}
	// 写入备份
	err = api.doc2File(docBackUpFilePath, docData)
	if err != nil {
		return err
	}

	// 更新 theme
	_ = new(APITheme).Set2File(themeObj)

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

func (api *APIDoc) doc2File(docFilePath string, docData []byte) error {
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
	return err
}

func (api *APIDoc) Get(theme string, docId string) (*entity.Doc, error) {
	data := &entity.Doc{}

	themeObj, err := new(APITheme).ThemeCacheGet(theme)
	if err != nil {
		logger.Error("not theme")
		return data, fmt.Errorf("not theme")
	}

	hashId, err := core.NewMetaDataObj.Get(docId)
	if err != nil {
		logger.Error("获取 theme Consistent hash id 失败 ", err)
	}
	logger.Info("hashId = ", hashId)
	meta, ok := themeObj.MetaData[hashId]
	if !ok {
		logger.Error("未找到元数据存储位置")
		return data, fmt.Errorf("未找到元数据存储位置")
	}

	docFilePath := meta.Routers + "/" + docId
	content, err := os.ReadFile(docFilePath)
	if err != nil {
		logger.Error("read file error:%v\n", err)
		// 去备份取
		content, err = os.ReadFile(meta.BackUp + "/" + docId)
		err = utils.DataDecoder(content, &data)
		if err != nil {
			return data, fmt.Errorf("no such doc.")
		}

		// TODO 备份变主数据，然后再备份到其他物理地址

	}
	err = utils.DataDecoder(content, &data)
	return data, err
}
