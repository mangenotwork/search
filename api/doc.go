package api

import (
	"bufio"
	"bytes"
	"compress/gzip"
	jsoniter "github.com/json-iterator/go"
	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils/logger"
	"io"
	"os"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type APIDoc struct {
}

// Set 写入文档
func (api *APIDoc) Set(doc *entity.Doc) {
	// 存储文档， 存在则更新
	docFilePath := "./doc/" + doc.ID
	docData, err := DataEncoder(doc)
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
}

func (api *APIDoc) Get(docId string) (*entity.Doc, error) {
	data := &entity.Doc{}
	docFilePath := "./doc/" + docId
	content, err := os.ReadFile(docFilePath)
	if err != nil {
		logger.Error("read file error:%v\n", err)
		return data, err
	}
	err = DataDecoder(content, &data)
	return data, err
}

// DataEncoder 数据量大，使用 json 序列化+gzip压缩
func DataEncoder(obj interface{}) ([]byte, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return []byte(""), err
	}
	return GzipCompress(b), nil
}

// DataDecoder 解码
func DataDecoder(data []byte, obj interface{}) error {
	b, err := GzipDecompress(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

// GzipCompress gzip压缩
func GzipCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := gzip.NewWriterLevel(&in, gzip.BestCompression)
	//w := gzip.NewWriter(&in)
	_, _ = w.Write(src)
	_ = w.Close()
	return in.Bytes()
}

// GzipDecompress gzip解压
func GzipDecompress(src []byte) ([]byte, error) {
	reader := bytes.NewReader(src)
	gr, err := gzip.NewReader(reader)
	if err != nil {
		return []byte(""), err
	}
	bf := make([]byte, 0)
	buf := bytes.NewBuffer(bf)
	_, err = io.Copy(buf, gr)
	err = gr.Close()
	return buf.Bytes(), nil
}
