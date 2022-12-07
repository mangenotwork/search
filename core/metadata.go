package core

import (
	"fmt"
	"github.com/mangenotwork/search/utils"
)

/*
虚拟分片单元数据

获取数据   一致性hash
  |
  V
虚拟层   master 记录虚拟层
  |
  V
寻址 数据真实存放位置
  |
  V
物理层 多个物理机集群

*/

var NewMetaDataObj *utils.Consistent

func NewMetaData() {
	NewMetaDataObj = utils.NewConsistent()
	// TODO 创建1万个虚拟分片
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("%d", i)
		NewMetaDataObj.Add(id)
	}
}
