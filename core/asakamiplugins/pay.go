package asakamiplugins

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	// PluginName 插件名字
	PluginVersion = "0.3.0"
	Help          = "稻荷宝帮助:\n- #稻荷宝 [pay/支付] [u[uid]/@] 金额\n- #稻荷宝 [balance/余额] \n- #稻荷宝 [help/帮助] \n- #稻荷宝 [金额] [jpy/dhc/dyc] [exchange/兑换] [jpy/dhc/dyc] \n- #稻荷宝 [rate/汇率] "
)

var (
// payfilepath = engine.DataFolder() + "pay/"
)

type paystruct struct {
	pluginversion string
	pluginhelp    string
	run           func(ctx *zero.Ctx)
}

var payv = paystruct{
	pluginversion: PluginVersion,
	pluginhelp:    Help,
	run:           payrun,
}

// Getpay 导出变量
func Getpay() paystruct {
	return payv
}

func payrun(ctx *zero.Ctx) {

	//log.Printf(ctx.Event.Message[0].Type)
	//log.Printf(ctx.Event.Message[1].Type)
	gid := ctx.Event.GroupID
	//uid := ctx.Event.UserID
	uid := link.Getrealid(ctx)
	itemdhc, err := backpack.GetItem(1001, uid)
	if err != nil {
		ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
		return
	}
	dhc := itemdhc.Quantity
	itemjpy, err := backpack.GetItem(1003, uid)
	if err != nil {
		ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
		return
	}
	jpy := itemjpy.Quantity
	itemdyc, err := backpack.GetItem(1002, uid)
	if err != nil {
		ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
		return
	}
	dyc := itemdyc.Quantity

	var msg = ctx.Event.Message
	text := msg.ExtractPlainText()

	if strings.Contains(text, "pay") || strings.Contains(text, "支付") {

		param := ctx.State["args"].(string)
		var fiancee int64
		var err error
		fianceer := regexp.MustCompile(`u\d+`)
		if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
			uidStr := ctx.Event.Message[1].Data["qq"]
			//log.Println(uidStr)
			fiancee, err = strconv.ParseInt(uidStr, 10, 64)
		} else if param == "" {
			fianceep := regexp.MustCompile(`^u`)
			fiancee, err = strconv.ParseInt(fianceep.ReplaceAllString(fianceer.FindString(text), ""), 10, 64)
		}

		//log.Println(fiancee)

		if err != nil {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 格式错误"))
			return
		}
		rtext := fianceer.ReplaceAllString(text, "")
		numr := regexp.MustCompile(`\d+\.?\d*`)
		num64, err := strconv.ParseFloat(numr.FindString(rtext), 64)
		num := float64(num64)

		if err != nil {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 格式错误"))
			return
		}
		res := ctx.GetGroupMemberInfo(gid, fiancee, true)

		if res.Type == gjson.Null {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 群员不存在"))
			return
		}

		if float64(dhc)-num < 0 {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 你的余额不足"))
			return
		}

		if uid == fiancee {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 不能给自己支付"))
			return
		}
		if strings.Contains(text, "jpy") || strings.Contains(text, "円") {
			if jpy-num < 0 {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 您的余额不足"))
				return
			}
			itemjpy.Quantity -= num
			backpack.UpdateItem(itemjpy, uid)
			fianceeitemjpy, err := backpack.GetItem(1003, fiancee)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			fianceeitemjpy.Quantity += num
			backpack.UpdateItem(fianceeitemjpy, fiancee)
			ctx.SendGroupMessage(gid, message.Text(fmt.Sprintf("支付成功！您当前的余额为%.02f円", jpy-num)))
			return
		} else if strings.Contains(text, "dhc") || strings.Contains(text, "稻荷币") {
			if dhc-num < 0 {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 您的余额不足"))
				return
			}
			itemdhc.Quantity -= num
			backpack.UpdateItem(itemdhc, uid)
			fianceeitemdhc, err := backpack.GetItem(1001, fiancee)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			fianceeitemdhc.Quantity += num
			backpack.UpdateItem(fianceeitemdhc, fiancee)

			ctx.SendGroupMessage(gid, message.Text(fmt.Sprintf("支付成功！您当前的余额为%.02f稻荷币", dhc-num)))
			return
		} else if strings.Contains(text, "dyc") || strings.Contains(text, "钓鱼币") {
			if dyc-num < 0 {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 您的余额不足"))
				return
			}
			itemdyc.Quantity -= num
			backpack.UpdateItem(itemdyc, uid)
			fianceeitemdyc, err := backpack.GetItem(1002, fiancee)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			fianceeitemdyc.Quantity += num
			backpack.UpdateItem(fianceeitemdyc, fiancee)

			ctx.SendGroupMessage(gid, message.Text(fmt.Sprintf("支付成功！您当前的余额为%.02f钓鱼币", dyc-num)))
			return
		}

	}
	if strings.Contains(text, "balance") || strings.Contains(text, "余额") {
		ctx.SendChain(message.At(uid), message.Text(fmt.Sprintf("余额：\n- %-06.02f\t稻荷币\n- %-06.02f\t円\n- %-06.02f\t钓鱼币", dhc, jpy, dyc)))
		return
	}
	if strings.Contains(text, "help") || strings.Contains(text, "帮助") {
		ctx.SendGroupMessage(gid, message.Text(Help))
		return
	}

	if strings.Contains(text, "exchange") || strings.Contains(text, "兑换") {

		numr := regexp.MustCompile(`\d+\.?\d*`)
		numstr := numr.FindString(text)
		num1, err := strconv.ParseFloat(numstr, 64)
		if !(strings.Contains(text, "all") || strings.Contains(text, "全部")) {
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 格式错误 01"))
				return
			}
		}

		//获取货币类型
		var item1 BackpackItem
		var item2 BackpackItem
		//在“兑换”前面的货币
		//用正则表达式匹配货币类型

		currencyr := regexp.MustCompile(`(jpy|円|dhc|稻荷币|dyc|钓鱼币)`)
		currencys := currencyr.FindAllString(text, -1)
		//log.Println(currencys)
		if len(currencys) != 2 {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 格式错误 02"))
			return
		}
		var currency1 CURRENCY
		var currency2 CURRENCY
		if currencys[0] == "jpy" || currencys[0] == "円" {
			item1, err = backpack.GetItem(1003, uid)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			currency1 = GetCurrency(1003)
		} else if currencys[0] == "dhc" || currencys[0] == "稻荷币" {
			item1, err = backpack.GetItem(1001, uid)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			currency1 = GetCurrency(1001)
		} else if currencys[0] == "dyc" || currencys[0] == "钓鱼币" {
			item1, err = backpack.GetItem(1002, uid)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			currency1 = GetCurrency(1002)
		} else {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 格式错误 03"))
			return
		}
		//在“兑换”后面的货币
		if currencys[1] == "jpy" || currencys[1] == "円" {
			item2, err = backpack.GetItem(1003, uid)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			currency2 = GetCurrency(1003)
		} else if currencys[1] == "dhc" || currencys[1] == "稻荷币" {
			item2, err = backpack.GetItem(1001, uid)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			currency2 = GetCurrency(1001)
		} else if currencys[1] == "dyc" || currencys[1] == "钓鱼币" {
			item2, err = backpack.GetItem(1002, uid)
			if err != nil {
				ctx.SendGroupMessage(gid, message.Text("ERROR: 未知错误"))
				return
			}
			currency2 = GetCurrency(1002)
		} else {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 格式错误 04"))
			return
		}

		//不能兑换相同的货币
		if item1.itemid == item2.itemid {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 不能兑换相同的货币 "))
			return
		}

		//如果兑换全部
		if strings.Contains(text, "all") || strings.Contains(text, "全部") {
			num1 = item1.Quantity
		}

		//判断余额是否足够
		if item1.Quantity-num1 < 0 {
			ctx.SendGroupMessage(gid, message.Text("ERROR: 您的余额不足"))
			return
		}

		//计算汇率
		rate := currency1.value / currency2.value

		//计算兑换结果
		item1.Quantity -= num1
		backpack.UpdateItem(item1, uid)
		item2.Quantity += num1 * rate
		backpack.UpdateItem(item2, uid)

		//兑换影响货币价值
		currency1.value -= currency1.value * (rand.Float64() * 0.0001)
		currency2.value += currency2.value * (rand.Float64() * 0.0001)

		//更新货币价值
		UpdateCurrency(currency1)
		UpdateCurrency(currency2)

		//输出结果

		ctx.SendGroupMessage(gid, message.Text(fmt.Sprintf("余额：\n- %-06.02f\t%s\n- %-06.02f\t%s", item1.Quantity, currency1.Topic, item2.Quantity, currency2.Topic)))
		return
	}

	if strings.Contains(text, "rate") || strings.Contains(text, "汇率") {
		//获取货币价值
		dhccurrency := GetCurrency(1001)
		jpycurrency := GetCurrency(1003)
		dyccurrency := GetCurrency(1002)
		//输出结果
		ctx.SendGroupMessage(gid, message.Text(fmt.Sprintf("价值：\n- 稻荷币 %-06.04f\n- 円 %-06.04f\n- 钓鱼币 %-06.04f", dhccurrency.value, jpycurrency.value, dyccurrency.value)))
		return
	}
	ctx.SendGroupMessage(gid, message.Text("ERROR: 未知命令 请使用 #稻荷宝 help 查看帮助"))

}

func init() {

}
