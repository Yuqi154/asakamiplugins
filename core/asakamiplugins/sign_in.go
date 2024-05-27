// Package score 签到
package asakamiplugins

import (
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/AnimeAPI/bilibili"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/imgfactory"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/golang/freetype"
	"github.com/wcharczuk/go-chart/v2"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	backgroundURL       = "https://iw233.cn/api.php?sort=pc"
	referer             = "https://weibo.com/"
	signinMax           = 1
	signinpluginversion = "0.4.0"
	// SCOREMAX 分数上限定为1200
	SCOREMAX   = 1200
	signinhelp = "- 签到\n- 获得签到背景[@xxx] | 获得签到背景\n- 设置签到预设(0~3)\n- 查看等级排名\n注:为跨群排名\n- 查看我的钱包\n- 查看钱包排名\n注:为本群排行，若群人数太多不建议使用该功能!!!"
)

var (
	cachePath = engine.DataFolder() + "cache/"
	rankArray = [...]int{0, 10, 20, 50, 100, 200, 350, 550, 750, 1000, 1200}
	styles    = []scoredrawer{
		drawScore17,
		drawScore16,
		drawScore15,
		drawScore17b2,
	}
)

type signinstruct struct {
	pluginversion string
	pluginhelp    string
	run1          func(ctx *zero.Ctx)
	run2          func(ctx *zero.Ctx)
	run3          func(ctx *zero.Ctx)
	run4          func(ctx *zero.Ctx)
	run5          func(ctx *zero.Ctx)
}

var signinv = signinstruct{
	pluginversion: signinpluginversion,
	pluginhelp:    signinhelp,
	run1:          signinrun1,
	run2:          signinrun2,
	run3:          signinrun3,
	run4:          signinrun4,
	run5:          signinrun5,
}

// Getsignin 导出变量
func Getsignin() signinstruct {
	return signinv
}

func signinrun1(ctx *zero.Ctx) {
	// 选择key
	key := ctx.State["regex_matched"].([]string)[1]
	gid := ctx.Event.GroupID
	if gid < 0 {
		// 个人用户设为负数
		gid = -link.Getrealid(ctx)
	}
	k := uint8(0)
	if key == "" {
		k = uint8(ctx.State["manager"].(*ctrl.Control[*zero.Ctx]).GetData(gid))
	} else {
		kn, err := strconv.Atoi(key)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		k = uint8(kn)
	}
	if int(k) >= len(styles) {
		ctx.SendChain(message.Text("ERROR: 未找到签到设定: ", key))
		return
	}
	uid := link.Getrealid(ctx)
	today := time.Now().Format("20060102")
	yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
	// 签到图片
	drawedFile := cachePath + strconv.FormatInt(uid, 10) + today + "signin.png"
	picFile := cachePath + strconv.FormatInt(uid, 10) + today + ".png"
	// 获取签到时间
	si := sdb.GetSignInByUID(uid)
	siUpdateTimeStr := si.UpdatedAt.Format("20060102")
	var csignin BackpackItem
	nowruntimeitem, err := backpack.GetItem(103, 0)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
	}

	switch {
	case si.Count >= signinMax && siUpdateTimeStr == today:
		// 如果签到时间是今天
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天你已经签到过了！"))
		if file.IsExist(drawedFile) {
			ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
		}
		return
	case siUpdateTimeStr != today && siUpdateTimeStr != yesterday && nowruntimeitem.Quantity >= 60*60*24:
		csignin, err := statistics.continuoussignin(uid)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		csignin.Quantity = 0
		backpack.UpdateItem(csignin, uid)
	case siUpdateTimeStr != today && siUpdateTimeStr == yesterday && nowruntimeitem.Quantity >= 60*60*24:
		// 如果是跨天签到就清数据
		err := sdb.InsertOrUpdateSignInCountByUID(uid, 0)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
	case siUpdateTimeStr == yesterday || nowruntimeitem.Quantity < 60*60*24:
		// 如果是连续签到
		var err error
		csignin, err = statistics.continuoussignin(uid)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}

	}
	// 更新签到次数
	err = sdb.InsertOrUpdateSignInCountByUID(uid, si.Count+1)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}

	//检查是否为2024年
	if time.Now().Year() == 2024 {
		statistics.signin(uid)
	}

	if csignin.Quantity == 0 {
		csignin, err = statistics.continuoussignin(uid)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
	}
	//log.Print(csignin.Quantity)

	// 更新经验
	level := sdb.GetScoreByUID(uid).Score + 1
	if level > SCOREMAX {
		level = SCOREMAX
		ctx.SendChain(message.At(uid), message.Text("你的等级已经达到上限"))
	}
	err = sdb.InsertOrUpdateScoreByUID(uid, level)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	// 更新钱包
	rank := getrank(level)
	dhbadd := 1.0 + rand.Float64()*10 + float64(rank*5) + csignin.Quantity*2 // 等级越高获得的钱越高
	jpyadd := rand.Float64()*float64(rank) + csignin.Quantity*0.1            // 等级越高获得的钱越高

	dhbitem, err := backpack.GetItem(1001, uid)
	if err != nil {
		ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
		return
	}
	renitem, err := backpack.GetItem(1003, uid)
	if err != nil {
		ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
		return
	}
	dhbitem.Quantity += dhbadd
	renitem.Quantity += jpyadd
	go func() {
		backpack.UpdateItem(dhbitem, uid)
		backpack.UpdateItem(renitem, uid)
	}()
	alldata := &scdata{
		drawedfile:       drawedFile,
		picfile:          picFile,
		uid:              uid,
		nickname:         ctx.CardOrNickName(uid),
		continuoussignin: int(csignin.Quantity),
		dhbinc:           dhbadd,
		jpyinc:           jpyadd,
		level:            level,
		rank:             rank,
	}
	drawimage, err := styles[k](alldata)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	// done.
	f, err := os.Create(drawedFile)
	if err != nil {
		data, err := imgfactory.ToBytes(drawimage)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.ImageBytes(data))
		return
	}
	_, err = imgfactory.WriteTo(drawimage, f)
	defer f.Close()
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
}

func signinrun2(ctx *zero.Ctx) {
	param := ctx.State["args"].(string)
	var uidStr string
	if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
		uidStr = ctx.Event.Message[1].Data["qq"]
	} else if param == "" {
		uidStr = strconv.FormatInt(link.Getrealid(ctx), 10)
	}
	picFile := cachePath + uidStr + time.Now().Format("20060102") + ".png"
	if file.IsNotExist(picFile) {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请先签到！"))
		return
	}
	if id := ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + picFile)); id.ID() == 0 {
		ctx.SendChain(message.Text("ERROR: 消息发送失败, 账号可能被风控"))
	}
}

func signinrun3(ctx *zero.Ctx) {

	today := time.Now().Format("20060102")
	drawedFile := cachePath + today + "scoreRank.png"
	if file.IsExist(drawedFile) {
		ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
		return
	}
	st, err := sdb.GetScoreRankByTopN(10)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	if len(st) == 0 {
		ctx.SendChain(message.Text("ERROR: 目前还没有人签到过"))
		return
	}
	_, err = file.GetLazyData(text.FontFile, control.Md5File, true)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	b, err := os.ReadFile(text.FontFile)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	font, err := freetype.ParseFont(b)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	f, err := os.Create(drawedFile)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	var bars []chart.Value
	for _, v := range st {
		if v.Score != 0 {
			bars = append(bars, chart.Value{
				Label: ctx.CardOrNickName(v.UID),
				Value: float64(v.Score),
			})
		}
	}
	err = chart.BarChart{
		Font:  font,
		Title: "等级排名(1天只刷新1次)",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: math.Ceil(bars[0].Value/10) * 10,
			},
		},
		Height:   500,
		BarWidth: 50,
		Bars:     bars,
	}.Render(chart.PNG, f)
	_ = f.Close()
	if err != nil {
		_ = os.Remove(drawedFile)
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
}

func signinrun4(ctx *zero.Ctx) {

	key := ctx.State["regex_matched"].([]string)[1]
	kn, err := strconv.Atoi(key)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	k := uint8(kn)
	if int(k) >= len(styles) {
		ctx.SendChain(message.Text("ERROR: 未找到签到设定: ", key))
		return
	}
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -link.Getrealid(ctx)
	}
	err = ctx.State["manager"].(*ctrl.Control[*zero.Ctx]).SetData(gid, int64(k))
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	ctx.SendChain(message.Text("设置成功"))
}

func init() {
	go func() {
		ok := file.IsExist(cachePath)
		if !ok {
			err := os.MkdirAll(cachePath, 0777)
			if err != nil {
				panic(err)
			}
			return
		}
		files, err := os.ReadDir(cachePath)
		if err == nil {
			for _, f := range files {
				if !strings.Contains(f.Name(), time.Now().Format("20060102")) {
					_ = os.Remove(cachePath + f.Name())
				}
			}
		}
	}()
	sdb = initialize(engine.DataFolder() + "db/level.db")
}

func getHourWord(t time.Time) string {
	h := t.Hour()
	switch {
	case 6 <= h && h < 12:
		return "早上好"
	case 12 <= h && h < 14:
		return "中午好"
	case 14 <= h && h < 19:
		return "下午好"
	case 19 <= h && h < 24:
		return "晚上好"
	case 0 <= h && h < 6:
		return "凌晨好"
	default:
		return ""
	}
}

func getrank(count int) int {
	for k, v := range rankArray {
		if count == v {
			return k
		} else if count < v {
			return k - 1
		}
	}
	return -1
}

func initPic(picFile string, uid int64) (avatar []byte, err error) {
	defer process.SleepAbout1sTo2s()
	avatar, err = web.GetData("http://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(uid, 10) + "&s=640")
	if err != nil {
		return
	}
	if file.IsExist(picFile) {
		return
	}
	url, err := bilibili.GetRealURL(backgroundURL)
	if err != nil {
		return
	}
	data, err := web.RequestDataWith(web.NewDefaultClient(), url, "", referer, "", nil)
	if err != nil {
		return
	}
	return avatar, os.WriteFile(picFile, data, 0644)
}

func signinrun5(ctx *zero.Ctx) {
	uid := link.Getrealid(ctx)
	signin2024, err := backpack.GetItem(102000, uid)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	continuoussignin, err := backpack.GetItem(102001, uid)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	ctx.SendChain(message.Text("2024年签到次数: ", signin2024.Quantity, "\n连续签到次数: ", continuoussignin.Quantity))
}
