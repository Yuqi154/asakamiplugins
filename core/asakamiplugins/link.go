package asakamiplugins

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type linkstruct struct {
	pluginversion  string
	pluginhelp     string
	run1           func(ctx *zero.Ctx)
	run2           func(ctx *zero.Ctx)
	Getrealid      func(ctx *zero.Ctx) int64
	Getrealidbyuid func(uid int64) int64
}

var linkv = linkstruct{
	pluginversion: "0.0.1",
	pluginhelp:    "link",
	run1: func(ctx *zero.Ctx) { //QQ
		uid := ctx.Event.UserID                         //获取用户KOOK ID
		revtext := ctx.Event.Message.ExtractPlainText() //获取消息
		fianceer := regexp.MustCompile(`\d+`)
		uidStr := fianceer.FindAllString(revtext, -1)
		if len(uidStr) == 0 {
			ctx.SendChain(message.Text("未找到QQ ID"))
			return
		}
		fiancee, err := strconv.ParseInt(uidStr[0], 10, 64)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		if fiancee == uid {
			ctx.SendChain(message.Text("不能绑定自己"))
			return
		}

		backpack.init(uid)
		var bond BackpackItem
		bond, err = backpack.GetItem(101, int64(uid))
		if err != nil {
			backpack.InsertItem(BackpackItem{ItemType: "SYSTEM", SubType: "BOND", itemid: 101, Quantity: 0}, int64(uid))
			bond, err = backpack.GetItem(101, int64(uid))
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
		}
		if bond.Quantity == 0 {
			//KOOK端确认绑定
			ctx.SendChain(message.Text("你确定要绑定QQ", fiancee, "吗？\n回答\"是\"或\"否\""))
			recv, cancel1 := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^(是|否)$`), zero.CheckUser(uid)).Repeat()
			defer cancel1()
			var confirm = false
			for {
				select {
				case <-time.After(time.Second * 60):
					ctx.SendChain(message.Text(ctx.Event.MessageID, message.Text("等待超时,取消绑定")))
					return
				case e := <-recv:
					nextcmd := e.Event.Message.String()
					if nextcmd == "否" {
						ctx.SendChain(message.Text(ctx.Event.MessageID, message.Text("已取消绑定")))
						return
					}
					confirm = true
				}
				if confirm {
					break
				}
			}
			//QQ端确认绑定
			code := strconv.Itoa(rand.Intn(8999) + 1000)
			ctx.SendChain(message.Text("请以您的QQ", fiancee, "向机器人发送验证码 ", code, " 以完成绑定"))
			recv, cancel2 := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^`+code+`$`), zero.CheckUser(fiancee)).Repeat()
			defer cancel2()
			confirm = false
			for {
				select {
				case <-time.After(time.Second * 60):
					ctx.SendChain(message.Text("等待超时,取消绑定"))
					return
				case <-recv:
					confirm = true
				}
				if confirm {
					break
				}
			}
			//绑定
			bond.Quantity = float64(fiancee)
			backpack.UpdateItem(bond, int64(uid))

			ctx.SendChain(message.Text("绑定成功"))
		}

		ctx.SendChain(message.Text("QQ ID: ", fiancee, "\n本地 ID: ", uid))
	},
	run2: func(ctx *zero.Ctx) { //KOOK
		//uid := ctx.Event.UserID                         //获取用户QQ ID
		revtext := ctx.Event.Message.ExtractPlainText() //获取消息
		fianceer := regexp.MustCompile(`\d+`)
		fiancee := fianceer.FindAllString(revtext, -1)

		if len(fiancee) == 0 {
			ctx.SendChain(message.Text("未找到KOOK ID"))
			return
		}
		ctx.SendChain(message.Text("暂不支持从QQ绑定KOOK ID"))
	},
	Getrealid: Getrealid,
}

func Getrealid(ctx *zero.Ctx) int64 {
	uid := ctx.Event.UserID //获取用户ID
	return Getrealidbyuid(uid)
}

func Getrealidbyuid(uid int64) int64 {
	bond, err := backpack.GetItem(101, int64(uid))
	if err != nil {
		backpack.InsertItem(BackpackItem{ItemType: "SYSTEM", SubType: "BOND", itemid: 101, Quantity: 0}, int64(uid))
		return uid
	}
	if bond.Quantity == 0 {
		return uid
	}
	fiancee := int64(bond.Quantity)
	return fiancee
}

func Getlink() linkstruct {
	return linkv
}
