package asakamiplugins

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type statisticss struct {
	pluginversion string
	pluginhelp    string
	run           func(ctx *zero.Ctx)
}

func Getstatistics() statisticss {
	return statisticss{
		pluginversion: "0.1.0",
		pluginhelp:    "- 浅上统计插件",
		run:           run,
	}
}

func run(ctx *zero.Ctx) {
	ctx.SendChain(message.Text("浅上统计插件"))
}

func (s statisticss) signin(uid int64) (BackpackItem, error) {
	signinv, err := backpack.GetItem(102000, uid)
	if err != nil {
		backpack.InsertItem(BackpackItem{
			ItemType: "SYSTEM",
			SubType:  "STATISTICS",
			itemid:   102000,
			Topic:    "签到次数",
			Tooltip:  "2024年签到次数",
			Buff:     "",
			Quantity: 1,
			Tag:      "SIGNIN",
		}, uid)
	} else {
		signinv.Quantity++
		backpack.UpdateItem(signinv, uid)
	}
	return signinv, nil
}

func (s statisticss) continuoussignin(uid int64) (BackpackItem, error) {
	continuoussigninv, err := backpack.GetItem(102001, uid)
	if err != nil {
		backpack.InsertItem(BackpackItem{
			ItemType: "SYSTEM",
			SubType:  "STATISTICS",
			itemid:   102001,
			Topic:    "连签次数",
			Tooltip:  "连续签到次数",
			Buff:     "",
			Quantity: 1,
			Tag:      "SIGNIN",
		}, uid)
	} else {
		continuoussigninv.Quantity++
		backpack.UpdateItem(continuoussigninv, uid)
	}
	return continuoussigninv, nil
}

// 计数总正常运行时间
func (s statisticss) sumruntime(time float64) (BackpackItem, error) {
	runtimev, err := backpack.GetItem(102, 0)
	if err != nil {
		backpack.InsertItem(BackpackItem{
			ItemType: "SYSTEM",
			SubType:  "STATISTICS",
			itemid:   102,
			Topic:    "总正常运行时间",
			Tooltip:  "总正常运行时间",
			Buff:     "",
			Quantity: time,
			Tag:      "RUNTIME",
		}, 0)
	} else {
		runtimev.Quantity++
		backpack.UpdateItem(runtimev, 0)
	}
	return runtimev, nil
}

// 计数本次正常运行时间
func (s statisticss) nowruntime(time float64) (BackpackItem, error) {
	runtimev, err := backpack.GetItem(103, 0)
	if err != nil {
		backpack.InsertItem(BackpackItem{
			ItemType: "SYSTEM",
			SubType:  "STATISTICS",
			itemid:   103,
			Topic:    "本次正常运行时间",
			Tooltip:  "本次正常运行时间",
			Buff:     "",
			Quantity: time,
			Tag:      "RUNTIME",
		}, 0)
	} else {
		runtimev.Quantity++
		backpack.UpdateItem(runtimev, 0)
	}
	return runtimev, nil
}

// 计数总离线时间
func (s statisticss) sumofflinetime(time float64) (BackpackItem, error) {
	offlinetimev, err := backpack.GetItem(104, 0)
	if err != nil {
		backpack.InsertItem(BackpackItem{
			ItemType: "SYSTEM",
			SubType:  "STATISTICS",
			itemid:   104,
			Topic:    "总离线时间",
			Tooltip:  "总离线时间",
			Buff:     "",
			Quantity: time,
			Tag:      "OFFLINETIME",
		}, 0)
	} else {
		offlinetimev.Quantity += time
		backpack.UpdateItem(offlinetimev, 0)
	}
	return offlinetimev, nil
}

// 计数上次离线时间
func (s statisticss) lastofflinetime(time float64) (BackpackItem, error) {
	offlinetimev, err := backpack.GetItem(105, 0)
	if err != nil {
		backpack.InsertItem(BackpackItem{
			ItemType: "SYSTEM",
			SubType:  "STATISTICS",
			itemid:   105,
			Topic:    "上次离线时间",
			Tooltip:  "上次离线时间",
			Buff:     "",
			Quantity: time,
			Tag:      "OFFLINETIME",
		}, 0)
	} else {
		offlinetimev.Quantity = time
		backpack.UpdateItem(offlinetimev, 0)
	}
	return offlinetimev, nil
}
