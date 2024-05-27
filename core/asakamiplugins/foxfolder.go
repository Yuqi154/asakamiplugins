package asakamiplugins

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"

	"os"

	"github.com/FloatTech/floatbox/file"
	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type foxfolderstruct struct {
	pluginversion string
	pluginhelp    string
	run           func(ctx *zero.Ctx)
	debug         func(ctx *zero.Ctx)
}

var foxfolderv = foxfolderstruct{
	pluginversion: "0.1.0",
	pluginhelp:    "- 叠狐狸 \n -#叠狐狸加入\n -#叠狐狸分析\n -#叠狐狸列表",
	run: func(ctx *zero.Ctx) {
		gid := ctx.Event.GroupID

		checkdb(gid)

		uid := link.Getrealid(ctx)
		weight, err := getweight(uid, gid)
		if err != nil {
			err = initweight(uid, gid)
			if err != nil {
				log.Print(err)
			}
		}
		if weight < 0.1 {
			weight = 0.1
			setweight(uid, weight, gid)
		}

		userlist, err := getuserlist(gid)
		if err != nil {
			log.Print(err)
		}
		weightlist, err := getweightlist(userlist, gid)
		if err != nil {
			log.Print(err)
		}

		length := len(userlist)

		rtext := ctx.Event.Message.ExtractPlainText()

		var sendtext string

		if strings.Contains(rtext, "加入") {
			//如果userlist里没有这个人，那么就加入
			for _, auid := range userlist {
				if auid == uid {
					sendtext += "你已经在狐狸堆里了"
					ctx.SendGroupMessage(ctx.Event.GroupID, message.ReplyWithMessage(ctx.Event.MessageID, message.Text(sendtext)))
					return
				}
			}

			leftlist, brokenlist := calprobability(userlist, uid, gid)
			if brokenlist != nil {
				sendtext += "有狐狸被压坏了 杂鱼~杂鱼~\n"
				sendtext += "你是狐狸堆中的第" + fmt.Sprint(len(leftlist)) + "只狐狸\n"
				for _, brokenid := range brokenlist {
					r := ctx.GetGroupMemberInfo(ctx.Event.GroupID, brokenid, true)
					//解析json
					var m map[string]interface{}
					err := json.Unmarshal([]byte(r.Raw), &m)
					if err != nil {
						log.Print(err)
					}
					//获取昵称
					nickname := m["nickname"].(string)
					sendtext += "♥ " + nickname + " ♥ " + fmt.Sprint(brokenid) + fmt.Sprintf(" %.02f", unwarpgetweight(brokenid, gid)) + "kg " + getdiscribe(unwarpgetweight(brokenid, gid)) + " 被压坏了\n"
				}
			} else {
				for _, auid := range leftlist {
					if auid == uid {
						sendtext += "加入成功，你是狐狸堆中的第" + fmt.Sprint(length+1) + "只狐狸"
						err = updateuserlist(leftlist, gid)
						if err != nil {
							log.Print(err)
						}
						ctx.SendGroupMessage(ctx.Event.GroupID, message.ReplyWithMessage(ctx.Event.MessageID, message.Text(sendtext)))
						return
					}
				}
				sendtext += "摔下来了啦 杂鱼~杂鱼~\n"
				weight -= rand.Float64()
				if weight < 0.1 {
					weight = 0.1
				}
				setweight(uid, weight, gid)
				sendtext += "你当前的体重是" + fmt.Sprintf("%.02f", weight) + "kg\n"
			}
			err = updateuserlist(leftlist, gid)
			if err != nil {
				log.Print(err)
			}
		} else if strings.Contains(rtext, "分析") {
			var totalweight float64
			for _, uid := range userlist {
				weight, err := getweight(uid, gid)
				if err != nil {
					log.Print(err)
				}
				totalweight += weight
			}
			if length == 0 {
				sendtext += "狐狸堆中没有狐狸\n"
			} else {
				sendtext += "狐狸堆中共有" + fmt.Sprint(length) + "只狐狸\n"
				sendtext += "狐狸堆的总重量是" + fmt.Sprintf("%.02f", totalweight) + "kg\n"
				sendtext += "狐狸堆的平均重量是" + fmt.Sprintf("%.02f", totalweight/float64(length)) + "kg\n"
			}
			successpossibility := executesuccesspossibility(int64(length), weight)

			successpossibility = successpossibility * 100
			sendtext += "你的体重是" + fmt.Sprintf("%.02f", weight) + "kg\n"
			var isin = false
			for _, auid := range userlist {
				if auid == uid {
					sendtext += "你已经在狐狸堆里了\n"
					isin = true
				}
			}
			if !isin {
				sendtext += "成功概率是" + fmt.Sprintf("%.02f", successpossibility) + "%\n"
				sendtext += "摔下概率是" + fmt.Sprintf("%.02f", 100-successpossibility) + "%\n"

				pic := drawPieChart(successpossibility / 100)

				ctx.SendGroupMessage(ctx.Event.GroupID, message.ReplyWithMessage(ctx.Event.MessageID, message.Text(sendtext)))
				ctx.SendGroupMessage(ctx.Event.GroupID, pic)
				return
			}
		} else if strings.Contains(rtext, "列表") {
			var nicknames []string
			//log.Print(userlist)
			for _, uid := range userlist {
				r := ctx.GetGroupMemberInfo(ctx.Event.GroupID, uid, true)
				//解析json
				var m map[string]interface{}
				err := json.Unmarshal([]byte(r.Raw), &m)
				if err != nil {
					log.Print(err)
				}
				//获取昵称
				nickname := m["nickname"].(string)
				nicknames = append(nicknames, nickname)
			}

			//如果列表里人数低于2人，那么就不绘制图表
			if length < 2 {
				if length == 1 {
					sendtext += "狐狸堆中只有1只狐狸\n"
					sendtext += "♥ " + nicknames[0] + " ♥ " + fmt.Sprint(userlist[0]) + " " + fmt.Sprintf("%.02f", weightlist[0]) + "kg " + getdiscribe(weightlist[0])
				}
				if length == 0 {
					sendtext += "狐狸堆中没有狐狸"
				}
				ctx.SendGroupMessage(ctx.Event.GroupID, message.Text(sendtext))
				return
			} else {
				sendtext += "狐狸堆中共有" + fmt.Sprint(length) + "只狐狸\n"
				for i := 0; i < length; i++ {
					sendtext += "♥ " + nicknames[i] + " ♥ " + fmt.Sprint(userlist[i]) + " " + fmt.Sprintf("%.02f", weightlist[i]) + "kg " + getdiscribe(weightlist[i]) + "\n"
				}

			}
			//绘制柱状图
			pic := drawBarChart(nicknames, weightlist, gid)
			//message.Message{pic, message.Text(sendtext)}
			//ctx.SendGroupMessage(ctx.Event.GroupID, message.Text(sendtext))
			ctx.SendGroupMessage(ctx.Event.GroupID, message.Message{message.Text(sendtext), pic})
			return
		}

		//清除最后一个换行
		sendtext = strings.TrimRight(sendtext, "\n")
		//ctx.SendGroupMessage(ctx.Event.GroupID, message.Text(sendtext))
		ctx.SendGroupMessage(ctx.Event.GroupID, message.ReplyWithMessage(ctx.Event.MessageID, message.Text(sendtext)))
	},
	debug: func(ctx *zero.Ctx) {
		if ctx.Event.GroupID != 1123489751 {
			return
		}
		rtext := ctx.Event.Message.ExtractPlainText()
		if strings.Contains(rtext, "清空") {
			cleanpile(ctx.Event.GroupID)
		}
	},
}

func drawPieChart(success float64) message.MessageSegment {

	fontfile, err := os.ReadFile("C:/Users/Administrator/AppData/Local/Microsoft/Windows/Fonts/HYZhongHeiTiS.ttf")
	if err != nil {
		log.Print(err)
	}
	font, err := truetype.Parse(fontfile)
	if err != nil {
		log.Print(err)
	}

	// Create a new pie chart
	pieChart := chart.PieChart{
		Font:   font,
		Width:  360,
		Height: 360,
		Values: []chart.Value{
			{Value: success, Label: "成功"},
			{Value: 1 - success, Label: "失败"},
		},
	}

	// Create a file to save the chart
	file, err := os.Create("data/asakamiplugins/foxfolder/pie_chart.png")
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	// Render the chart to the file
	err = pieChart.Render(chart.PNG, file)
	if err != nil {
		log.Print(err)
	}

	path, _ := os.Getwd()
	path = strings.ReplaceAll(path, "\\", "/")
	// Send the chart image as a message
	return message.Image("file:///" + path + "/data/asakamiplugins/foxfolder/pie_chart.png")
}

func drawBarChart(nicknames []string, weights []float64, gid int64) message.MessageSegment {
	// Create a new bar chart
	length := len(nicknames)
	fontfile, err := os.ReadFile("C:/Users/Administrator/AppData/Local/Microsoft/Windows/Fonts/HYZhongHeiTiS.ttf")
	if err != nil {
		log.Print(err)
	}
	font, err := truetype.Parse(fontfile)
	if err != nil {
		log.Print(err)
	}
	//barss = []chart.Value{}
	barChart := chart.BarChart{
		Title: "狐狸堆",
		TitleStyle: chart.Style{
			FontSize:  24,
			Font:      font,
			FillColor: chart.ColorWhite,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 80,
			},
			FillColor: chart.ColorWhite,
		},
		Height: 360,
		//Width:    720,
		BarWidth: 32,
		//BarSpacing:   32,
		Bars:         []chart.Value{{Value: 0, Label: "0"}},
		Font:         font,
		UseBaseValue: false,
		BaseValue:    0,
	}

	//if length > 5 {
	//barChart.Width = 720 + (length-5)*64
	//}

	// Add the nicknames as bars to the chart
	for i := 0; i < length; i++ {
		barChart.Bars = append(barChart.Bars, chart.Value{Value: weights[i], Label: nicknames[i], Style: chart.Style{TextRotationDegrees: 180}})
	}

	//删去第一个空的bar
	barChart.Bars = barChart.Bars[1:]

	// Create a file to save the chart
	file, err := os.Create("data/asakamiplugins/foxfolder/bar_chart_" + fmt.Sprint(gid) + ".png") // "bar_chart.png
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	// Render the chart to the file
	err = barChart.Render(chart.PNG, file)
	if err != nil {
		log.Print(err)
	}

	path, _ := os.Getwd()
	path = strings.ReplaceAll(path, "\\", "/")
	// Send the chart image as a message
	return message.Image("file:///" + path + "/data/asakamiplugins/foxfolder/bar_chart_" + fmt.Sprint(gid) + ".png")
}

// Getfoxfolder 导出变量
func Getfoxfolder() foxfolderstruct {
	return foxfolderv
}

func getdiscribe(weight float64) string {
	var discribe string
	if weight < 1 {
		discribe = "奶狐"
	} else if weight < 5 {
		discribe = "小狐狸"
	} else if weight < 10 {
		discribe = "大狐狸"
	} else if weight < 20 {
		discribe = "幼年狐娘"
	} else {
		discribe = "成年狐娘"
	}
	return discribe
}

func createTables(db *sql.DB, gid int64) {
	// Create the first table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS weight_` + fmt.Sprint(gid) + ` (
		uid INT64,
		weight FLOAT64
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the second table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pile_` + fmt.Sprint(gid) + ` (
		serial INT64,
		uid INT64
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化个人重量
func initweight(uid int64, gid int64) error {
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/foxfolder.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO weight_"+fmt.Sprint(gid)+" (uid, weight) VALUES (?, ?)", uid, 0.1)
	if err != nil {
		return err
	}
	return nil
}

func getweight(uid int64, gid int64) (float64, error) {
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/foxfolder.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var weight float64
	err = db.QueryRow("SELECT weight FROM weight_"+fmt.Sprint(gid)+" WHERE uid=?", uid).Scan(&weight)
	if err != nil {
		initweight(uid, gid)
	}
	return weight, nil
}

func unwarpgetweight(uid int64, gid int64) float64 {
	a, e := getweight(uid, gid)
	if e != nil {
		log.Print(e)
	}
	return a
}

func setweight(uid int64, weight float64, gid int64) error {
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/foxfolder.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE weight_"+fmt.Sprint(gid)+" SET weight=? WHERE uid=?", weight, uid)
	if err != nil {
		return err
	}
	return nil
}

// 获取群组内所有用户列表
func getuserlist(gid int64) ([]int64, error) {
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/foxfolder.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var uidlist []int64
	rows, err := db.Query("SELECT uid FROM pile_" + fmt.Sprint(gid))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var uid int64
		err = rows.Scan(&uid)
		if err != nil {
			return nil, err
		}
		uidlist = append(uidlist, uid)
	}
	return uidlist, nil
}

func updateuserlist(uidlist []int64, gid int64) error {

	//log.Print(uidlist)
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/foxfolder.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM pile_" + fmt.Sprint(gid))
	if err != nil {
		return err
	}
	//先清空表
	db.Exec("DELETE FROM sqlite_sequence WHERE name='pile_" + fmt.Sprint(gid) + "'")
	//重置自增序列
	db.Exec("VACUUM")
	//清理数据库
	//再插入
	for _, uid := range uidlist {
		_, err = db.Exec("INSERT INTO pile_"+fmt.Sprint(gid)+" (uid) VALUES (?)", uid)
		if err != nil {
			return err
		}
	}
	return nil
}

func getweightlist(uidlist []int64, gid int64) ([]float64, error) {
	var weightlist []float64
	for _, uid := range uidlist {
		weight, err := getweight(uid, gid)
		if err != nil {
			return nil, err
		}
		weightlist = append(weightlist, weight)
	}
	return weightlist, nil
}

// 计算加入群组的概率 成功，失败，压坏,返回剩余的id和摔下的id
func calprobability(uidlist []int64, uid int64, gid int64) ([]int64, []int64) {
	weightlist, err := getweightlist(uidlist, gid)
	if err != nil {
		log.Print(err)
	}

	//计算堆叠重量
	var totalweight float64
	for _, weight := range weightlist {
		totalweight += weight
	}

	uidweight, err := getweight(uid, gid)
	if err != nil {
		log.Print(err)
	}

	//计算是否成功
	var success bool

	successpossibility := executesuccesspossibility(int64(len(uidlist)), uidweight)
	//判定是否成功，rand.Float64()是(0,1)的均匀分布
	if rand.Float64() < successpossibility {
		success = true
	} else {
		success = false
	}

	if success {
		uidlist = append(uidlist, uid)
		weightlist = append(weightlist, uidweight)
		//如果列表里只有自己，那么就不需要计算
		if len(uidlist) == 1 {
			return uidlist, nil
		}

		//检查每个人身上的重量，创建一个重量列表
		var perweightlist []float64
		for i := 0; i < len(weightlist); i++ {
			var perweight float64
			for j := i + 1; j < len(weightlist); j++ {
				perweight += weightlist[j]
			}
			perweightlist = append(perweightlist, perweight)
		}
		//如果一个人身上的重量是她自身的重量的十倍以上，那么她有0.5的概率就会被压坏
		var brokenid int64
		var brokenidlist []int64
		var leftidlist []int64
		for i := 0; i < len(perweightlist)-1; i++ {
			if perweightlist[i] > weightlist[i]*10 && rand.Float64() < 0.5 {
				brokenid = uidlist[i]
				pweight, err := getweight(brokenid, gid)
				if err != nil {
					log.Print(err)
				}
				pweight *= rand.Float64()*0.2 + 0.8
				setweight(brokenid, pweight, gid)
				brokenidlist = append(brokenidlist, brokenid)
			} else if perweightlist[i] > weightlist[i] {
				//有概率压坏=承重/(体重*10)/2
				if rand.Float64() < perweightlist[i]/(weightlist[i]*10)/2 {
					brokenid = uidlist[i]
					pweight, err := getweight(brokenid, gid)
					if err != nil {
						log.Print(err)
					}
					pweight *= rand.Float64()*0.1 + 0.9
					setweight(brokenid, pweight, gid)
					brokenidlist = append(brokenidlist, brokenid)
				} else {
					leftidlist = append(leftidlist, uidlist[i])
				}
			} else {
				leftidlist = append(leftidlist, uidlist[i])
			}

		}
		leftidlist = append(leftidlist, uid)
		uidweight += float64(len(brokenidlist)) * rand.Float64()
		setweight(uid, uidweight, gid)
		return leftidlist, brokenidlist
	} else {
		return uidlist, nil
	}
}

func executesuccesspossibility(length int64, uidweight float64) float64 {
	var successpossibility float64

	successpossibility = 5 * math.Log(float64(length)+3) / (float64(length) + 3) / 1.5

	if uidweight < 20 {
		successpossibility *= 1 - (math.Sqrt((uidweight+5)/20) - 0.5)
	} else {
		successpossibility *= 1 - (math.Sqrt(1.25) - 0.5)
	}

	if successpossibility > 1 {
		successpossibility = 1
	}
	if successpossibility < 0 {
		successpossibility = 0
	}

	return successpossibility
}

func cleanpile(gid int64) {
	//清空堆内所有狐狸
	updateuserlist(nil, gid)
}

func init() {
	//创建数据库文件夹
	if file.IsNotExist("data/asakamiplugins/foxfolder") {
		err := os.MkdirAll("data/asakamiplugins/foxfolder", 0755)
		if err != nil {
			panic(err)
		}
	}

	// Connect to the database
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/foxfolder.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}

func checkdb(gid int64) {
	// Connect to the database
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/foxfolder.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//检查表是否存在
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='weight_" + fmt.Sprint(gid) + "'").Scan(&name)
	if err != nil {
		log.Print(err)
	}
	if strings.Contains(name, "weight_"+fmt.Sprint(gid)) {
		return
	} else {
		createTables(db, gid)
	}
}
