package utils

import (
	"bytes"
	"compress/gzip"
	jsoniter "github.com/json-iterator/go"
	"io"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
