package asakamiplugins

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	rollPluginVersion = "0.1.0"
	rollHelp          = "- 冬香酱写的roll功能\n  使用方法：roll 【标签】 骰子数量d骰子面数(!)/判定点数\n  Note:骰子面数可以为4,6,8,10,12,20,100\n  例子：roll tag 1d6/4\n  也可以跳过面数检查：roll tag 1d16!/4"
)

var (
	mg = []string{"早安", "早", "お早う", "一键三连"}
	ng = []string{"晚安", "晚", "Good night", "Sweet dreams", "Sleep tight", "Rest well", "Nighty-night", "おやすみなさい", "いい夢を", "ぐっすり眠って", "安らかに", "お休みなさい", "一键三连"}
)

type rollstruct struct {
	pluginversion string
	pluginhelp    string
	run           func(ctx *zero.Ctx)
}

var rollv = rollstruct{
	pluginversion: rollPluginVersion,
	pluginhelp:    rollHelp,
	run:           rollrun,
}

type sleepmanage struct {
	pluginversion string
	pluginhelp    string
	run1          func(ctx *zero.Ctx)
	run2          func(ctx *zero.Ctx)
}

var sleepmanagev = sleepmanage{
	pluginversion: "0.1.0",
	pluginhelp:    "",
	run1:          sleepmanagerun1,
	run2:          sleepmanagerun2,
}

type chatstruct struct {
	pluginversion string
	pluginhelp    string
	sleepmanage   sleepmanage
	roll          rollstruct
}

var chatv = chatstruct{
	pluginversion: "0.1.0",
	pluginhelp:    "chat插件帮助:\n" + rollHelp,
	sleepmanage:   sleepmanagev,
	roll:          rollv,
}

func Getchat() chatstruct {
	return chatv
}

func rollrun(ctx *zero.Ctx) {
	msg_s := ctx.Event.Message.ExtractPlainText()
	gid := ctx.Event.GroupID

	tag := getindex(msg_s, "^roll ", " ")
	msg_p := regexp.MustCompile("^roll ")
	msg := msg_p.ReplaceAllString(msg_s, "")
	num_s := getindex(msg, " ", "/")
	log.Println(num_s)
	point, errp := strconv.Atoi(getindex(msg, "/", ""))

	num := 0
	face := 6
	if strings.Contains(num_s, "d") {
		errn := error(nil)
		errf := error(nil)
		num, errn = strconv.Atoi(getindex(msg, " ", "d"))
		if strings.Contains(num_s, "!") {
			face, errf = strconv.Atoi(getindex(msg, "d", "\\!"))
			if errf != nil {
				ctx.SendGroupMessage(gid, message.Text("["+tag+"] 格式错误 E06"))
				return
			}
		} else {
			face, errf = strconv.Atoi(getindex(msg, "d", "/"))
		}
		log.Print(num, face)
		log.Print(errn, errf)
		if errn != nil || errf != nil {
			ctx.SendGroupMessage(gid, message.Text("["+tag+"] 格式错误 E01"))
			return
		}

		if face != 4 && face != 6 && face != 8 && face != 10 && face != 12 && face != 20 && face != 100 {
			if strings.Contains(num_s, "!") {
				face, errf = strconv.Atoi(getindex(msg, "d", "\\!"))
				if errf != nil {
					ctx.SendGroupMessage(gid, message.Text("["+tag+"] 格式错误 E02"))
					return
				}
			} else {
				ctx.SendGroupMessage(gid, message.Text("["+tag+"] 骰子面数不合法 E05"))
				return
			}
		}

	} else {
		errn := error(nil)
		num, errn = strconv.Atoi(num_s)
		if errn != nil {
			ctx.SendGroupMessage(gid, message.Text("["+tag+"] 格式错误 E03"))
			return
		}
	}

	if errp != nil {
		ctx.SendGroupMessage(gid, message.Text("["+tag+"] 格式错误 E04"))
		return
	}
	//随机数
	sum := 0
	for ; num > 0; num-- {
		randnum := rand.Int63n(int64(face))
		sum += int(randnum) + 1
	}
	//判定成功
	sus := ""
	if sum < point {
		sus = " 成功"
	} else {
		sus = " 失败"
	}

	ctx.SendGroupMessage(gid, message.Text("["+tag+"] "+strconv.Itoa(sum)+"~"+strconv.Itoa(point)+sus))
}

func sleepmanagerun1(ctx *zero.Ctx) {
	position, getUpTime := ssdb.getUp(ctx.Event.GroupID, link.Getrealid(ctx))
	log.Debugln(position, getUpTime)
	hour, minute, second := timeDuration(getUpTime)
	if (hour == 0 && minute == 0 && second == 0) || hour >= 24 {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("早安成功！你是今天第%d个起床的", position)))
	} else {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("早安成功！你的睡眠时长为%d时%d分%d秒,你是今天第%d个起床的", hour, minute, second, position)))
	}
}

func sleepmanagerun2(ctx *zero.Ctx) {
	position, sleepTime := ssdb.sleep(ctx.Event.GroupID, link.Getrealid(ctx))
	log.Debugln(position, sleepTime)
	hour, minute, second := timeDuration(sleepTime)
	if (hour == 0 && minute == 0 && second == 0) || hour >= 24 {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("晚安成功！你是今天第%d个睡觉的", position)))
	} else {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("晚安成功！你的清醒时长为%d时%d分%d秒,你是今天第%d个睡觉的", hour, minute, second, position)))
	}
}

func init() {
	go func() {
		ssdb = sinitialize(engine.DataFolder() + "db/manage.db")
	}()
}

func timeDuration(time time.Duration) (hour, minute, second int64) {
	hour = int64(time) / (1000 * 1000 * 1000 * 60 * 60)
	minute = (int64(time) - hour*(1000*1000*1000*60*60)) / (1000 * 1000 * 1000 * 60)
	second = (int64(time) - hour*(1000*1000*1000*60*60) - minute*(1000*1000*1000*60)) / (1000 * 1000 * 1000)
	return hour, minute, second
}

// 只统计6点到12点的早安
func IsMorning(*zero.Ctx) bool {
	now := time.Now().Hour()
	return now >= 6 && now <= 12
}

// 只统计21点到凌晨3点的晚安
func IsEvening(*zero.Ctx) bool {
	now := time.Now().Hour()
	return now >= 21 || now <= 3
}
