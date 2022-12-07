package api

import (
	"bufio"
	"fmt"
	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils"
	"github.com/mangenotwork/search/utils/logger"
	"os"
	"path/filepath"
)

type APITheme struct {
}

var ThemeCache map[string]*entity.Theme

// ThemeCacheInit ThemeCache 初始化
func (api *APITheme) ThemeCacheInit() {
	ThemeCache = map[string]*entity.Theme{}
	for _, v := range api.GetAll() {
		logger.Info("初始化 Theme = ", v)
		ThemeCache[v.Name] = v
	}
}

func (api *APITheme) ThemeCacheGet(theme string) (*entity.Theme, error) {
	v, ok := ThemeCache[theme]
	if !ok {
		return nil, fmt.Errorf("nil")
	}
	return v, nil
}

func (api *APITheme) Created(theme *entity.Theme) error {
	themeFilePath := entity.ThemePath + theme.Name
	logger.Info("theme file = ", themeFilePath)
	if utils.FileExists(themeFilePath) {
		logger.Error("已经创建")
		return fmt.Errorf("theme存在")
	}

	themeData, err := utils.DataEncoder(theme)
	if err != nil {
		logger.Error("文档不能被压缩, err = ", err)
		return err
	}
	file, err := os.OpenFile(themeFilePath, os.O_WRONLY|os.O_CREATE, 0666)
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
	_, err = write.Write(themeData)
	if err != nil {
		logger.Error("文件写入失败, err = ", err)
	}
	// Flush将缓存的文件真正写入到文件中
	err = write.Flush()
	if err != nil {
		logger.Error("写入到文件中失败, err = ", err)
	}

	// 更新缓存
	api.ThemeCacheInit()

	return nil
}

func (api *APITheme) Set2File(theme *entity.Theme) error {
	themeFilePath := entity.ThemePath + theme.Name
	themeData, err := utils.DataEncoder(theme)
	if err != nil {
		logger.Error("文档不能被压缩, err = ", err)
		return err
	}
	file, err := os.OpenFile(themeFilePath, os.O_WRONLY|os.O_CREATE, 0666)
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
	_, err = write.Write(themeData)
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
	return nil
}

func (api *APITheme) get(path string) (*entity.Theme, error) {
	data := &entity.Theme{}
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("read file error:%v\n", err)
		return data, fmt.Errorf("not this theme.")
	}
	err = utils.DataDecoder(content, &data)
	return data, err
}

func (api *APITheme) Get(name string) (*entity.Theme, error) {
	themeFilePath := entity.ThemePath + name
	return api.get(themeFilePath)
}

func (api *APITheme) GetAll() []*entity.Theme {
	data := make([]*entity.Theme, 0)
	filepath.Walk(entity.ThemePath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				fmt.Println("dir:", path)
				return nil
			}
			fmt.Println("file:", path, filepath.Base(path))
			d, err := api.get(path)
			if err != nil {
				logger.Error("获取Theme数据失败, err = ", err, " | path = ", path)
			}
			data = append(data, d)
			return nil
		})
	return data
}

func (api *APITheme) GetAllName() []string {
	data := make([]string, 0)
	return data
}

// TODO 文档总大小
// TODO 索引总大小
// TODO 重置索引
// TODO 清空索引
// TODO 文档导出 - excel
// TODO 文档导出 - mysql
