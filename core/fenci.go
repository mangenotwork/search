package core

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mangenotwork/search/entity"
	"github.com/mangenotwork/search/utils"
	"github.com/mangenotwork/search/utils/logger"

	"github.com/go-ego/gse"
)

var (
	seg gse.Segmenter
)

// 加载分词词典
func init() {
	// load default dict
	err := seg.LoadDict()
	//err := seg.LoadDictEmbed()
	if err != nil {
		panic("[致命错误] 加载分词词典失败！")
	}

	// 初始化 元数据分片对象
	NewMetaData()
}

// TermExtract 提取索引词
// 除了标点符号，助词，语气词，形容词，叹词 其他都被分出来
func TermExtract(str string) []*entity.Term {
	segments := seg.Segment([]byte(str))
	termList := make([]*entity.Term, 0)
	for _, v := range segments {
		t := v.Token()
		p := t.Pos()
		txt := t.Text()
		end := v.End()
		start := v.Start()
		//logger.Info("txt = ", txt, p)

		if p == "w" || p == "u" || p == "y" || p == "a" || p == "e" {
			continue
		}

		if p == "x" && !utils.ContainsEnglishAndNumber(txt) {
			continue
		}

		termList = append(termList, &entity.Term{
			Text:  txt,
			Freq:  t.Freq(),
			End:   end,
			Start: start,
		})
	}
	return termList
}

/*
GetTermCase

	保留词规则 : n + x(英文) + i
	名词	n	n
	名语素	ng	ng
	人名	nr	nr
	地名	ns	ns
	机构团体	nt
	外文字符	nx
	其他专名	nz
	i  成语
	j  简略语
*/
func GetTermCase(str string) []string {
	segments := seg.Segment([]byte(str))
	termList := make([]string, 0)
	for _, v := range segments {
		t := v.Token()
		p := t.Pos()
		txt := t.Text()
		logger.Info("txt = ", txt, p)
		if strings.Index(p, "n") == -1 && p != "x" && p != "i" && p != "j" {
			continue
		}
		if p == "x" && !utils.ContainsEnglish(txt) {
			continue
		}
		termList = append(termList, txt)
	}
	return termList
}

// PLInfo 索引信息文件
//
// .pli
// []*pliFile  file_name,valMax,valMin,
// fNum 文件数量
//
// 每个索引包含 一个 信息文件 主要记录索引存储结构,每个文件最多存储100条数据
// 这样设计的缺点: 空间浪费，写慢
// 这样设计的优点: 读快，读取结果已经被三个维度分别排序
//
// 文件：
// .plt 的文件 postingList time  按时间排序的数据存储  k:docId  v:time 排序规则:只有文档时间和词频两个维度 t>f
// .plo 的文件 postingList orderInt 按排序值来进行排序 k:docId  v:orderInt 排序规则:有自定义排序值，文档时间，词频三个维度 o>t>f
// .plf 的文件 postingList Freq 按词频值来进行排序  k:docId   v:Freq  排序规则: 只有文档时间和词频两个维度  f>t
//
// 结构:
// 存储结构  []*d{docId, value(用来排序的), start, end}
type PLInfo struct {
	PLDir   string        // 保存数据的路径
	PLTFile []*PLFileInfo // 起到一个游标的作用, 时间排序数据文件
	PLOFile []*PLFileInfo // 排序值排序数据文件
	PLFFile []*PLFileInfo // 按词频排序数据文件
	PLTFNum int           // 文件数量
	PLOFNum int
	PLFFNum int
	StartF  int // 启始数， 1开始， 是最大的排序值的数据
	EndF    int // 结束数，= FNum , 是最小的排序值的数据
}

type PLFileInfo struct {
	FileId   int     // 文件编号
	FileName string  // 文件名称
	ValMax   float64 // 最大排序值
	ValMin   float64 // 最小排序值
	Num      int     // 数据条数
}

// 倒排表写入到磁盘
func setPostings(docId, dir string, docStamp, orderInt float64, term *entity.Term) {
	plDir := dir + term.Text + "/"
	// 创建一个 term目录
	utils.Mkdir(plDir)
	// 是否存在信息文件，没有就初始化
	pliFile := plDir + "i.pli"
	if !utils.FileExists(pliFile) {
		setPlInfo(plDir, pliFile)
	}
	// 获取当前信息
	pliObj, err := getPlInfo(pliFile)
	if err != nil {
		logger.Error("获取词索引信息文件失败, err = ", err)
		return
	}
	setPLData(docId, "plt", docStamp, term.Freq, 0, pliObj, term)
	setPLData(docId, "plo", orderInt, docStamp, term.Freq, pliObj, term)
	setPLData(docId, "plf", term.Freq, docStamp, 0, pliObj, term)
}

func SetPostings(theme, docId, text string, docStamp, orderInt int64) {
	dir := entity.IndexPath + theme + "/"
	for _, v := range TermExtract(text) {
		setPostings(docId, dir, float64(docStamp), float64(orderInt), v)
	}
}

func SetPostingsAuthor(theme, docId, text string, docStamp, orderInt int64) {
	dir := entity.IndexPath + theme + "/"
	list := TermExtract(text)
	list = append(list, &entity.Term{
		Text:  text,
		Freq:  1,
		End:   0,
		Start: len(text),
	})
	for _, v := range list {
		setPostings(docId, dir, float64(docStamp), float64(orderInt), v)
	}
}

func setPLData(docId, plType string, value, sortPara1, sortPara2 float64, pliObj *PLInfo, term *entity.Term) {
	// 一条数据都不存在的情况
	firstFile := fmt.Sprintf("%s1.%s", pliObj.PLDir, plType)
	if !utils.FileExists(firstFile) {
		setData2FileInit(plType, value, pliObj, &entity.PL{docId, value, term.End, term.Start,
			sortPara1, sortPara2})
		return
	}

	// 只有一个文件，数据量小于100
	var plFileObj []*PLFileInfo
	num := 0
	switch plType {
	case "plt":
		plFileObj = pliObj.PLTFile
		num = pliObj.PLTFNum
	case "plo":
		plFileObj = pliObj.PLOFile
		num = pliObj.PLOFNum
	case "plf":
		plFileObj = pliObj.PLFFile
		num = pliObj.PLFFNum
	}
	//logger.Info("当前 num 与  plFileObj[0].Num = ", num, plFileObj[0].Num)
	if num == 1 && plFileObj[0].Num < 100 {
		// 直接读取文件写入
		setData2File(plType, docId, value, sortPara1, sortPara2, pliObj, term)
		return
	}

	// 定位文件， 通过排序值定位这条数据应该被插入到哪个文件
	notPos := false
	thisMax := false
	thisMin := false
	logger.Error("执行范围查找 plFileObj = ", plFileObj)
	for _, infoO := range plFileObj {
		// 在这个范围
		if value <= infoO.ValMax && value >= infoO.ValMin {
			notPos = true
			logger.Info("范围定在这个文件 ", infoO.FileName)
			// 读这个文件，然后写入数据，如果超过 100 条，将最后一个数据写入下一个文件  这个可以使用回调
			data2File(infoO.FileId, plType, docId, value, sortPara1, sortPara2, pliObj, term)
			break
		}
		if value > infoO.ValMax {
			thisMax = true
		}
		if value < infoO.ValMin {
			thisMin = true
		}
	}

	// 最大 写在第一个文件
	if !notPos && thisMax {
		data2File(1, plType, docId, value, sortPara1, sortPara2, pliObj, term)
		return
	}

	// 最小 写在第二个文件
	if !notPos && thisMin {
		data2File(pliObj.EndF, plType, docId, value, sortPara1, sortPara2, pliObj, term)
		return
	}
}

func data2File(fileNum int, plType, docId string, value, sortPara1, sortPara2 float64, pliObj *PLInfo, term *entity.Term) {
	plFile := fmt.Sprintf("%s%d.%s", pliObj.PLDir, fileNum, plType)
	// 不存在就创建 就结束递归
	if !utils.FileExists(plFile) {
		setData2FileNumInit(fileNum, plType, pliObj, &entity.PL{docId, value, term.End, term.Start,
			sortPara1, sortPara2})
		return
	}

	plList := getPlListFromFile(plFile)

	// 查看是否存在， 如何存在就替换，不存在就增加
	pltIsSet := false
	i := 0
	valueList := make([]float64, 0)
	for _, pltData := range plList {
		valueList = append(valueList, pltData.Value)

		// 存在更新
		if pltData.Key == docId {
			plList[i] = &entity.PL{docId, value, term.End, term.Start, sortPara1, sortPara2}
			valueList[i] = value
			pltIsSet = true
			goto Has
		}
		i++
	}

	// 不存在写入
	if !pltIsSet {
		plList = append(plList, &entity.PL{docId, value, term.End, term.Start, sortPara1, sortPara2})
		valueList = append(valueList, value)
		i++
	}

Has:
	// 排序
	sortRule(plType, plList)

	// 文件是否已经满100条数据
	var sendPlList []*entity.PL
	thisFull := false
	if plType == "plo" {
		logger.Error("数据条数 : ", len(plList), " | i = ", i)
	}
	if i >= 100 {
		sendPlList = plList[0:100]
		if plType == "plo" {
			logger.Error("超过100个数据 !!!!!!!!!! ")
		}
		thisFull = true
	} else {
		sendPlList = plList
	}

	// 存储
	pltListData, err := utils.DataEncoder(sendPlList)
	if err != nil {
		logger.Error("压缩 pltList 失败 : ", err)
	}
	setFile(plFile, pltListData)

	// 更新信息文件
	data2FileUpdatePLI(plType, plFile, len(plList), fileNum, pliObj, valueList)

	// 超过100 数据往后移 直到插满
	if thisFull {
		nextData := plList[len(plList)-1]
		logger.Info("超过100 数据往后移动直到插满 !!!! 数据 = ", nextData)
		fileNum++
		data2File(fileNum, plType, nextData.Key, nextData.Value, nextData.SortPara1, nextData.SortPara2, pliObj,
			&entity.Term{
				End:   nextData.End,
				Start: nextData.Start,
			})
	}
}

// data2File 的更新信息文件方法
func data2FileUpdatePLI(plType, plFile string, dataNum, fileNum int, pliObj *PLInfo, valueList []float64) {
	switch plType {
	case "plt":
		has := false
		for index, f := range pliObj.PLTFile {
			if f.FileName == plFile {
				has = true
				pliObj.PLTFile[index].ValMax = getMax(valueList)
				pliObj.PLTFile[index].ValMin = getMin(valueList)
				pliObj.PLTFile[index].Num = dataNum
				break
			}
		}
		if !has {
			pliObj.PLTFile = append(pliObj.PLTFile, &PLFileInfo{
				FileId:   fileNum + 1,
				FileName: plFile,
				ValMax:   getMax(valueList),
				ValMin:   getMin(valueList),
				Num:      dataNum,
			})
			pliObj.PLFFNum++
		}

	case "plo":
		has := false
		for index, f := range pliObj.PLOFile {
			if f.FileName == plFile {
				has = true
				pliObj.PLOFile[index].ValMax = getMax(valueList)
				pliObj.PLOFile[index].ValMin = getMin(valueList)
				pliObj.PLOFile[index].Num = dataNum
				break
			}
		}
		if !has {
			pliObj.PLOFile = append(pliObj.PLOFile, &PLFileInfo{
				FileId:   fileNum + 1,
				FileName: plFile,
				ValMax:   getMax(valueList),
				ValMin:   getMin(valueList),
				Num:      dataNum,
			})
			pliObj.PLOFNum++
		}

	case "plf":
		has := false
		for index, f := range pliObj.PLFFile {
			if f.FileName == plFile {
				has = true
				pliObj.PLFFile[index].ValMax = getMax(valueList)
				pliObj.PLFFile[index].ValMin = getMin(valueList)
				pliObj.PLFFile[index].Num = dataNum
				break
			}
		}
		if !has {
			pliObj.PLFFile = append(pliObj.PLFFile, &PLFileInfo{
				FileId:   fileNum + 1,
				FileName: plFile,
				ValMax:   getMax(valueList),
				ValMin:   getMin(valueList),
				Num:      dataNum,
			})
			pliObj.PLFFNum++
			pliObj.EndF = pliObj.PLFFNum
		}

	}
	updatePlInfo(pliObj)
}

// 排序 排序规则
// .plt 的文件 postingList time  按时间排序的数据存储  k:docId  v:time 排序规则:只有文档时间和词频两个维度 t>f
//
// .plo 的文件 postingList orderInt 按排序值来进行排序 k:docId  v:orderInt 排序规则:有自定义排序值，文档时间，词频三个维度 o>t>f
//
// .plf 的文件 postingList Freq 按词频值来进行排序  k:docId   v:Freq  排序规则: 只有文档时间和词频两个维度  f>t
func sortRule(plType string, plList []*entity.PL) {
	switch plType {
	case "plt":
		sort.Slice(plList, func(i, j int) bool {
			if plList[i].Value == plList[j].Value {
				return plList[i].SortPara1 > plList[j].SortPara1
			}
			return plList[i].Value > plList[j].Value
		})
	case "plo":
		sort.Slice(plList, func(i, j int) bool {
			if plList[i].Value == plList[j].Value {
				if plList[i].SortPara1 == plList[j].SortPara1 {
					return plList[i].SortPara2 > plList[j].SortPara2
				}
				return plList[i].SortPara1 > plList[j].SortPara1
			}
			return plList[i].Value > plList[j].Value
		})
	case "plf":
		sort.Slice(plList, func(i, j int) bool {
			if plList[i].Value == plList[j].Value {
				return plList[i].SortPara1 > plList[j].SortPara1
			}
			return plList[i].Value > plList[j].Value
		})
	}
}

func getPlListFromFile(plFile string) []*entity.PL {
	plList := make([]*entity.PL, 0)
	content, err := os.ReadFile(plFile)
	if err != nil {
		logger.Error("read file error: ", err)
		return plList
	}
	err = utils.DataDecoder(content, &plList)
	if err != nil {
		logger.Error("pltFile 读取数据失败 = ", err)
		return plList
	}
	return plList
}

// 直接写入文件
func setData2File(plType, docId string, value, sortPara1, sortPara2 float64, pliObj *PLInfo, term *entity.Term) {
	plFile := fmt.Sprintf("%s%d.%s", pliObj.PLDir, 1, plType)
	plList := getPlListFromFile(plFile)

	// 查看是否存在， 如何存在就替换，不存在就增加
	pltIsSet := false
	i := 0
	valueList := make([]float64, 0)
	for _, pltData := range plList {
		valueList = append(valueList, pltData.Value)

		// 存在更新
		if pltData.Key == docId {
			plList[i] = &entity.PL{docId, value, term.End, term.Start,
				sortPara1, sortPara2}
			valueList = append(valueList, value)
			pltIsSet = true
			break
		}
		i++
	}

	// 不存在写入
	if !pltIsSet {
		plList = append(plList, &entity.PL{docId, value, term.End, term.Start,
			sortPara1, sortPara2})
		valueList = append(valueList, value)
	}

	// 排序
	sortRule(plType, plList)

	// 存储
	pltListData, err := utils.DataEncoder(plList)
	if err != nil {
		logger.Error("压缩 pltList 失败 : ", err)
	}
	setFile(plFile, pltListData)

	// 更新信息文件
	switch plType {
	case "plt":
		pliObj.PLTFile[0].Num = len(plList)

	case "plo":
		pliObj.PLOFile[0].Num = len(plList)

	case "plf":
		pliObj.PLFFile[0].Num = len(plList)

	}
	updatePlInfo(pliObj)
}

// 将数据写到文件
func setFile(plFile string, data []byte) {
	logger.Info("写入文件  ======> ", plFile)
	file, err := os.OpenFile(plFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("文件打开失败 err = ", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()
	write := bufio.NewWriter(file)
	_, err = write.Write(data)
	if err != nil {
		logger.Error("写入文件失败, err = ", err)
	}
	err = write.Flush()
	if err != nil {
		logger.Error("写入文件失败, err = ", err)
	}
}

func setData2FileInit(plType string, value float64, pliObj *PLInfo, data *entity.PL) {
	plList := make([]*entity.PL, 0)
	plList = append(plList, data)
	// 存储
	pltListData, err := utils.DataEncoder(plList)
	if err != nil {
		logger.Error("压缩 pltList 失败 : ", err)
	}
	plFile := fmt.Sprintf("%s%d.%s", pliObj.PLDir, 1, plType)
	setFile(plFile, pltListData)

	// 更新信息 TODO 每种类型可以扩展
	switch plType {
	case "plt":
		pliObj.PLTFile = append(pliObj.PLTFile, &PLFileInfo{
			FileId:   1,
			FileName: plFile,
			ValMax:   value,
			ValMin:   value,
			Num:      1,
		})
	case "plo":
		pliObj.PLOFile = append(pliObj.PLOFile, &PLFileInfo{
			FileId:   1,
			FileName: plFile,
			ValMax:   value,
			ValMin:   value,
			Num:      1,
		})
	case "plf":
		pliObj.PLFFile = append(pliObj.PLFFile, &PLFileInfo{
			FileId:   1,
			FileName: plFile,
			ValMax:   value,
			ValMin:   value,
			Num:      1,
		})
	}
	updatePlInfo(pliObj)
}

// 不存在就新建
func setData2FileNumInit(num int, plType string, pliObj *PLInfo, data *entity.PL) {
	plList := []*entity.PL{data}
	pltListData, err := utils.DataEncoder(plList)
	if err != nil {
		logger.Error("压缩 pltList 失败 : ", err)
		return
	}
	plFile := fmt.Sprintf("%s%d.%s", pliObj.PLDir, num, plType)
	setFile(plFile, pltListData)

	// 更新信息 TODO 每种类型可扩展
	switch plType {
	case "plt":
		pliObj.PLTFile = append(pliObj.PLTFile, &PLFileInfo{
			FileId:   num,
			FileName: plFile,
			ValMax:   data.Value,
			ValMin:   data.Value,
			Num:      1,
		})
		pliObj.PLTFNum++

	case "plo":
		pliObj.PLOFile = append(pliObj.PLOFile, &PLFileInfo{
			FileId:   num,
			FileName: plFile,
			ValMax:   data.Value,
			ValMin:   data.Value,
			Num:      1,
		})
		pliObj.PLOFNum++

	case "plf":
		pliObj.PLFFile = append(pliObj.PLFFile, &PLFileInfo{
			FileId:   num,
			FileName: plFile,
			ValMax:   data.Value,
			ValMin:   data.Value,
			Num:      1,
		})
		pliObj.PLFFNum++
		pliObj.EndF = pliObj.PLFFNum

	}
	updatePlInfo(pliObj)
}

func setPlInfo(plDir, pliFile string) {
	plInfo := &PLInfo{
		PLDir:   plDir,
		PLTFNum: 1,
		PLOFNum: 1,
		PLFFNum: 1,
		StartF:  1,
		EndF:    1,
	}

	plInfo.PLTFile = make([]*PLFileInfo, 0)
	plInfo.PLOFile = make([]*PLFileInfo, 0)
	plInfo.PLFFile = make([]*PLFileInfo, 0)

	pltListData, err := utils.DataEncoder(plInfo)
	if err != nil {
		logger.Error("压缩 pltList 失败 : ", err)
	}
	setFile(pliFile, pltListData)
}

func getPlInfo(pliFile string) (*PLInfo, error) {
	plInfo := &PLInfo{}
	content, err := os.ReadFile(pliFile)
	if err != nil {
		logger.Errorf("read file error:%v\n", err)
		return nil, err
	}
	err = utils.DataDecoder(content, &plInfo)
	if err != nil {
		logger.Error("pltFile 读取数据失败 = ", err)
	}
	logger.Info("获取到的 getPlInfo = ", plInfo)
	return plInfo, err
}

func updatePlInfo(pliObj *PLInfo) {
	pliFile := pliObj.PLDir + "i.pli"
	plInfoData, err := utils.DataEncoder(pliObj)
	if err != nil {
		logger.Error("压缩 plInfoData 失败 : ", err)
	}
	setFile(pliFile, plInfoData)
}

func getMax(val []float64) float64 {
	var max float64
	for _, v := range val {
		if v > max {
			max = v
		}
	}
	return max
}

func getMin(val []float64) float64 {
	var min float64
	for _, v := range val {
		if min == 0 || v <= min {
			min = v
		}
	}
	return min
}
