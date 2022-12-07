/*
将niupp 的文章写入到
*/
package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
	"strings"
)

var (
	host192         = "192.168.4.2"
	port            = 3306
	user            = "root"
	password        = "admin"
	spiderBase      = "niu_pp"
	SpiderBaseDB, _ = gt.NewMysql(host192, port, user, password, spiderBase)
)

func main() {
	//data, _ := SpiderBaseDB.Query("select * from tbl_topic where fid=1 group by tid limit 105;")
	data, _ := SpiderBaseDB.Query("select * from tbl_topic where tid < 400000")
	for _, v := range data {
		//log.Println(v)

		/*
			{
			    "id": "1",
				"title": "阿伤悲三45645阿仕顿渤海湾怕萨芬请勿恐怕内插法",
				"author": "阿斯达",
				"time_stamp": 1669884252,
				"order_int": 0,
				"content": "综合俄新社、《消息报》等多家俄媒报道，欧盟委员会主席乌尔苏拉·冯德莱恩11月30日早些时候在推特上发了一段视频，她在视频中称，俄罗斯对乌克兰发起特别军事行动以来，超10万名乌克兰士兵和超2万平民已经死亡。随后，她用一段新视频替换了原视频，新视频中关于乌军损失数据的片段“消失了”。对此，俄罗斯联邦安全会议副主席梅德韦杰夫当天晚些时候作出回应。\n\n冯德莱恩 资料图冯德莱恩 资料图\n　　梅德韦杰夫当天在社交媒体VK上先是写道：“疯狂的乌尔苏拉大妈在自己的推特账号上向世界宣布，乌克兰武装部队蒙受‘10万损失’。随后，这一消息被删除了，经过了修改。”经典美文，听写——让心情舒畅！每天上午10点准时与大家相见。Hints:Southern Africarly African American: Jumping the Broom In the times of slavery in this country, African American couples were not allowed to formally marry and live together. To make a public declaration of their love and commitment, a man and woman jumped over a broom into matrimony, to the beat of drums. The broom has long held significant meaning for the various Africans, symbolizing, the start of home - making for the newlywed couple. In Southern Africa, the day after the wedding, the bride assisted the other women in the family in sweeping the courtyard, indicating her dutiful willing ness to help her in-laws with housework till the newlyweds could move to their new home. Some African-American couples today are choosing to include this symbolic rite in their wedding ceremony.早期非洲裔美国人：跳扫帚在美国的黑奴时代，黑人男女是不允许正式结婚生活在一起的。为了向世人宣布他们的.爱情和婚约，一对黑人男女和着鼓声的节奏，一起跳过一把扫帚。（扫帚对各种非洲人长期来都具有很重要的意义，因为它意味着新婚夫妇组成家庭的开始。在南部非洲，新娘在婚后的第一天要帮助夫家的其他女性清扫院子，以此表明在住进自己的新家前，她愿意尽职地帮助丈夫的家人承担家务劳动。）直至今日，一些美国黑人还在他们的婚礼上举行这种象征性的仪式。这篇材料你能听出多少？点击这里做听写，提高外语水平>>",
				"description": ""
			}
		*/

		game := ""
		if v["gameCode"] == "200" {
			game = "竞彩"
		}
		if v["gameCode"] == "1" {
			game = "福彩"
		}
		if v["gameCode"] == "2" {
			game = "双色球"
		}
		if v["gameCode"] == "4" {
			game = "排列三"
		}
		if v["gameCode"] == "5" {
			game = "大乐透"
		}

		j := `{
			    "id": "` + v["tid"] + `",
				"title": "` + game + `",
				"author": "` + v["nickName"] + `",
				"time_stamp": ` + v["created"] + `,
				"order_int": ` + v["orderInt"] + `,
				"content": "` + strings.Replace(v["content"], "\"", "\\\"", -1) + `",
				"description": ""
			}`
		log.Println(j)
		ctx, _ := gt.PostJson("http://127.0.0.1:14444/doc/case1", j)
		gt.Info(ctx.Json)
	}

}
