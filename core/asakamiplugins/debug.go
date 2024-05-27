package asakamiplugins

import (
	"log"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type debugstruct struct {
	pluginversion string
	pluginhelp    string
	run1          func(ctx *zero.Ctx)
	run2          func(ctx *zero.Ctx)
}

var debugv = debugstruct{
	pluginversion: "0.0.1",
	pluginhelp:    "debug",
	run1: func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		gid := ctx.Event.GroupID
		realuid := link.Getrealid(ctx)
		ctx.SendGroupMessage(gid, message.Text("uid:", uid, "\ngid:", gid, "\nrealuid:", realuid))
	},
	run2: func(ctx *zero.Ctx) {
		sumruntimeitem, err := backpack.GetItem(102, 0)
		if err != nil {
			log.Print(err)
		}
		nowruntimeitem, err := backpack.GetItem(103, 0)
		if err != nil {
			log.Print(err)
		}
		sumofflinetimeitem, err := backpack.GetItem(104, 0)
		if err != nil {
			log.Print(err)
		}
		lastofflinetimeitem, err := backpack.GetItem(105, 0)
		if err != nil {
			log.Print(err)
		}
		sumruntime := int64(sumruntimeitem.Quantity)
		nowruntime := int64(nowruntimeitem.Quantity)
		sumofflinetime := int64(sumofflinetimeitem.Quantity)
		lastofflinetime := int64(lastofflinetimeitem.Quantity)
		//将秒数转换为天时分秒
		sumofflinedays := int(sumofflinetime / 86400)
		sumofflinehours := int(sumofflinetime % 86400 / 3600)
		sumofflineminutes := int(sumofflinetime % 3600 / 60)
		sumofflineseconds := int(sumofflinetime % 60)
		lastofflinedays := int(lastofflinetime / 86400)
		lastofflinehours := int(lastofflinetime % 86400 / 3600)
		lastofflineminutes := int(lastofflinetime % 3600 / 60)
		lastofflineseconds := int(lastofflinetime % 60)
		nowrundays := int(nowruntime / 86400)
		nowrunhours := int(nowruntime % 86400 / 3600)
		nowrunminutes := int(nowruntime % 3600 / 60)
		nowrunseconds := int(nowruntime % 60)
		sumrundays := int(sumruntime / 86400)
		sumrunhours := int(sumruntime % 86400 / 3600)
		sumrunminutes := int(sumruntime % 3600 / 60)
		sumrunseconds := int(sumruntime % 60)
		ctx.SendChain(message.Text(sumofflinetimeitem.Topic, ":", sumofflinedays, "天", sumofflinehours, "时", sumofflineminutes, "分", sumofflineseconds, "秒\n", lastofflinetimeitem.Topic, ":", lastofflinedays, "天", lastofflinehours, "时", lastofflineminutes, "分", lastofflineseconds, "秒\n", nowruntimeitem.Topic, ":", nowrundays, "天", nowrunhours, "时", nowrunminutes, "分", nowrunseconds, "秒\n", sumruntimeitem.Topic, ":", sumrundays, "天", sumrunhours, "时", sumrunminutes, "分", sumrunseconds, "秒"))
	},
}

func Getdebug() debugstruct {
	return debugv
}
