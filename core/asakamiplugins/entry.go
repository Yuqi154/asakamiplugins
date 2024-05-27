package asakamiplugins

import (
	"os"

	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	cachePath := engine.DataFolder() + "cache/"
	go func() {
		_ = os.RemoveAll(cachePath)
		err := os.MkdirAll(cachePath, 0755)
		if err != nil {
			panic(err)
		}
	}()
	engine.OnPrefixGroup([]string{"#浅上", "#ASAKAMI"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		corerun(ctx)
	})

	//注册插件aipaint
	engine.OnPrefix("#春生画图").SetBlock(false).Handle(func(ctx *zero.Ctx) {
		aipaint.run(ctx)
	})

	//注册插件signin
	engine.OnRegex(`^签到\s?(\d*)$|^一键三连$|^签到.*$`).Limit(ctxext.LimitByUser).SetBlock(false).Handle(func(ctx *zero.Ctx) {
		signin.run1(ctx)
	})
	engine.OnPrefixGroup([]string{"获得签到背景", "获取签到背景"}, zero.OnlyGroup).Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		signin.run2(ctx)
	})
	engine.OnFullMatch("查看等级排名", zero.OnlyGroup).Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		signin.run3(ctx)
	})
	engine.OnRegex(`^设置签到预设\s*(\d+)$`, zero.SuperUserPermission).Limit(ctxext.LimitByUser).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		signin.run4(ctx)
	})
	engine.OnFullMatch(`签到数据`).Limit(ctxext.LimitByUser).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		signin.run4(ctx)
	})

	//注册插件pay
	engine.OnPrefix(`#稻荷宝`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		pay.run(ctx)
	})

	//注册插件chat
	engine.OnPrefix("roll").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		chat.roll.run(ctx)
	})
	engine.OnFullMatchGroup(mg, IsMorning, zero.OnlyGroup).SetBlock(false).Handle(func(ctx *zero.Ctx) {
		chat.sleepmanage.run1(ctx)
	})
	engine.OnFullMatchGroup(ng, IsEvening, zero.OnlyGroup).SetBlock(false).Handle(func(ctx *zero.Ctx) {
		chat.sleepmanage.run2(ctx)
	})

	//注册插件backpack
	engine.OnFullMatch("#初始化背包").SetBlock(false).Handle(func(ctx *zero.Ctx) {
		backpack.init(link.Getrealid(ctx))
		ctx.SendGroupMessage(ctx.Event.GroupID, message.Text("初始化背包成功"))
	})

	//debug
	engine.OnPrefix("#debug#uid").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		debug.run1(ctx)
	})
	engine.OnPrefix("#debug#runtime").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		debug.run2(ctx)
	})

	//link
	engine.OnPrefix("#绑定QQ").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		link.run1(ctx)
	})
	engine.OnPrefix("#绑定KOOK").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		link.run2(ctx)
	})

	//fox
	engine.OnPrefix("#叠狐狸").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		games.foxfolder.run(ctx)
	})
	engine.OnPrefix("#DEBUG叠狐狸").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		games.foxfolder.debug(ctx)
	})
	//gpt
	/*
		engine.OnPrefix("#GPT").SetBlock(false).Handle(func(ctx *zero.Ctx) {
			gpt.run(ctx)
		})
		engine.OnMessage(zero.OnlyToMe).SetBlock(false).Handle(func(ctx *zero.Ctx) {
			gpt.reply(ctx)
		})*/
	//rcon
	engine.OnPrefix("#rcon").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		rcon.run(ctx)
	})
}
