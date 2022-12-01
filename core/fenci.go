package core

import (
	"bufio"
	"fmt"
	"github.com/go-ego/gse"
	"github.com/mangenotwork/search/utils/logger"
	"os"
	"strings"
)

var (
	seg gse.Segmenter
)

func init() {
	// load default dict
	//err := seg.LoadDict()

	err := seg.LoadDictEmbed()
	if err != nil {
		panic("segment error")
	}
}

func Case() {
	segCut()
}

// TermExtract 提取索引词
// 除了标点符号，其他都被分出来
func TermExtract(str string) []string {
	segments := seg.Segment([]byte(str))
	termList := make([]string, 0)
	for _, v := range segments {
		t := v.Token()
		p := t.Pos()
		txt := t.Text()
		logger.Info("txt = ", txt, p)
		if p == "w" {
			continue
		}
		if p == "x" && !ContainsEnglishAndNumber(txt) {
			continue
		}

		termList = append(termList, txt)
	}
	return termList
}

/*
保留词规则 : n + x(英文) + i
//名词	n	n
//名语素	ng	ng
//人名	nr	nr
//地名	ns	ns
//机构团体	nt
//外文字符	nx
//其他专名	nz
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
		if p == "x" && !ContainsEnglish(txt) {
			continue
		}

		termList = append(termList, txt)
	}
	return termList
}

func segCut() {
	// 文本id
	docId := "2"

	// 分词文本
	tb := "综合俄新社、《消息报》等多家俄媒报道，欧盟委员会主席乌尔苏拉·冯德莱恩11月30日早些时候在推特上发了一段视频，她在视频中称，俄罗斯对乌克兰发起特别军事行动以来，超10万名乌克兰士兵和超2万平民已经死亡。随后，她用一段新视频替换了原视频，新视频中关于乌军损失数据的片段“消失了”。对此，俄罗斯联邦安全会议副主席梅德韦杰夫当天晚些时候作出回应。\n\n冯德莱恩 资料图冯德莱恩 资料图\n　　梅德韦杰夫当天在社交媒体VK上先是写道：“疯狂的乌尔苏拉大妈在自己的推特账号上向世界宣布，乌克兰武装部队蒙受‘10万损失’。随后，这一消息被删除了，经过了修改。”经典美文，听写——让心情舒畅！每天上午10点准时与大家相见。Hints:Southern Africarly African American: Jumping the Broom In the times of slavery in this country, African American couples were not allowed to formally marry and live together. To make a public declaration of their love and commitment, a man and woman jumped over a broom into matrimony, to the beat of drums. The broom has long held significant meaning for the various Africans, symbolizing, the start of home - making for the newlywed couple. In Southern Africa, the day after the wedding, the bride assisted the other women in the family in sweeping the courtyard, indicating her dutiful willing ness to help her in-laws with housework till the newlyweds could move to their new home. Some African-American couples today are choosing to include this symbolic rite in their wedding ceremony.早期非洲裔美国人：跳扫帚在美国的黑奴时代，黑人男女是不允许正式结婚生活在一起的。为了向世人宣布他们的.爱情和婚约，一对黑人男女和着鼓声的节奏，一起跳过一把扫帚。（扫帚对各种非洲人长期来都具有很重要的意义，因为它意味着新婚夫妇组成家庭的开始。在南部非洲，新娘在婚后的第一天要帮助夫家的其他女性清扫院子，以此表明在住进自己的新家前，她愿意尽职地帮助丈夫的家人承担家务劳动。）直至今日，一些美国黑人还在他们的婚礼上举行这种象征性的仪式。这篇材料你能听出多少？点击这里做听写，提高外语水平>>"

	segments := seg.Segment([]byte(tb))

	for _, v := range segments {
		t := v.Token()
		p := t.Pos()
		txt := t.Text()

		if strings.Index(p, "n") == -1 && p != "x" {
			continue
		}
		if p == "x" && !ContainsEnglish(txt) {
			continue
		}

		logger.Info(txt, " \t |", "词频 = ", v.Token().Freq(), " | 词性标注 = ", p)

		// 这个 txt 就是一个  term

		// 创建一个 term目录
		Mkdir("./index/" + txt)

		invertedFile := "./index/" + txt + "/" + "1.t"

		postingList := make([]string, 0)

		contentStr := ""
		if !isExist(invertedFile) {
			postingList = append(postingList, docId)
		} else {
			// 读取索引
			content, err := os.ReadFile(invertedFile)
			if err != nil {
				fmt.Printf("read file error:%v\n", err)
				return
			}
			contentStr = string(content)
			postingList = strings.Split(contentStr, ";")
			isSet := false
			for _, v := range postingList {
				if v == docId {
					isSet = true
					break
				}
			}
			if !isSet {
				postingList = append(postingList, docId)
			}
		}

		file, err := os.OpenFile(invertedFile, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {

			fmt.Println("文件打开失败", err)

		}

		//及时关闭file句柄
		defer file.Close()
		//写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)
		write.WriteString(strings.Join(postingList, ";"))
		// Flush将缓存的文件真正写入到文件中
		write.Flush()

	}

}

func ContainsEnglish(str string) bool {
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for _, v := range str {
		if strings.Contains(dictionary, string(v)) {
			return true
		}
	}
	return false
}

func ContainsEnglishAndNumber(str string) bool {
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for _, v := range str {
		if strings.Contains(dictionary, string(v)) {
			return true
		}
	}
	return false
}

// 目录是否存在，不存在则创建
func Mkdir(term string) {
	err := os.MkdirAll(term, os.ModePerm)
	if err != nil {
		return
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
