package asakamiplugins

import (
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	//将插件名字放入数组，方便后续调用
	engine = control.Register("asakamiplugins", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              help,
		PrivateDataFolder: "asakamiplugins",
	})
	coreversion = "0.3.2"
	help        = "- 浅上插件核心\n- 使用#浅上help获取更多帮助"

	//获取插件结构体
	aipaint     = Getaipaint()
	pay         = Getpay()
	signin      = Getsignin()
	chat        = Getchat()
	gpt         = Getgpt()
	economy     = Geteconomy()
	backpack    = Getbackpack()
	rpg         = Getrpg()
	statistics  = Getstatistics()
	debug       = Getdebug()
	link        = Getlink()
	games       = Getgames()
	rcon        = GetRcon()
	pluginarray = []string{"asakamiplugins", "pay", "signin", "chat", "gpt", "aipaint", "economy", "backpack", "rpg", "statistics", "debug", "link", "games"}
)

func corerun(ctx *zero.Ctx) {
	text := ctx.Event.Message.ExtractPlainText()
	//uid := ctx.Event.UserID
	gid := ctx.Event.GroupID
	if strings.Contains(text, "help") || strings.Contains(text, "帮助") {
		str := ""
		for _, plugin := range pluginarray {
			if strings.Contains(text, plugin) {
				str += plugin + ": " + GetPluginHelp(plugin) + "\n"
			}
		}
		if str == "" {
			ctx.SendGroupMessage(gid, message.Text("命令列表：\n#浅上 help\n#浅上 status\n#浅上 version\n#浅上 list\n#浅上 [插件名] help\n#浅上 [插件名] status\n#浅上 [插件名] version\n"))
		} else {
			ctx.SendGroupMessage(gid, message.Text(str))
		}
		return

	}
	if strings.Contains(text, "status") || strings.Contains(text, "状态") {
		//输出插件状态
		for _, plugin := range pluginarray {
			if strings.Contains(text, plugin) {
				ctx.SendGroupMessage(gid, message.Text(plugin+": "+GetPluginStatus(plugin)))
			}
		}
		return
	}
	if strings.Contains(text, "reload") || strings.Contains(text, "重载") {
		//占位
		return
	}
	if strings.Contains(text, "update") || strings.Contains(text, "更新") {
		//占位
		return
	}
	if strings.Contains(text, "version") || strings.Contains(text, "版本") {
		//输出核心版本
		ctx.SendGroupMessage(gid, message.Text("asakamicore"+coreversion))
		//再遍历一次数组，获得命令中的插件名字
		str := ""
		for _, plugin := range pluginarray {
			if strings.Contains(text, plugin) {
				str += plugin + ": " + GetPluginVersion(plugin) + "\n"
			}
		}
		ctx.SendGroupMessage(gid, message.Text(str))
		return
	}
	if strings.Contains(text, "list") || strings.Contains(text, "列表") {
		//输出插件列表
		str := ""
		for _, plugin := range pluginarray {
			str += plugin + ": " + GetPluginVersion(plugin) + "\n"
		}
		ctx.SendGroupMessage(gid, message.Text(str))
		return
	}

}

func GetPluginVersion(plugin string) (version string) {
	switch plugin {
	case "pay":
		version = pay.pluginversion
	case "signin":
		version = signin.pluginversion
	case "chat":
		version = chat.pluginversion
	case "plugins":
		version = coreversion
	case "gpt":
		version = gpt.pluginversion
	case "aipaint":
		version = aipaint.pluginversion
	case "economy":
		version = economy.pluginversion
	case "backpack":
		version = backpack.pluginversion
	case "rpg":
		version = rpg.pluginversion
	case "statistics":
		version = statistics.pluginversion
	case "debug":
		version = debug.pluginversion
	case "link":
		version = link.pluginversion
	case "games":
		version = games.pluginversion
	}
	return
}

func GetPluginStatus(plugin string) (status string) {
	service, ok := control.Lookup(plugin)
	if service == nil {
		status = "未找到插件"
	}
	if ok {
		status = "已启用"
	} else {
		status = "已禁用"
	}
	return status
}

func GetPluginHelp(plugin string) (version string) {
	switch plugin {
	case "pay":
		version = pay.pluginhelp
	case "signin":
		version = signin.pluginhelp
	case "chat":
		version = chat.pluginhelp
	case "asakamiplugins":
		version = coreversion
	case "gpt":
		version = gpt.pluginhelp
	case "aipaint":
		version = aipaint.pluginhelp
	case "economy":
		version = economy.pluginhelp
	case "backpack":
		version = backpack.pluginhelp
	case "rpg":
		version = rpg.pluginhelp
	case "statistics":
		version = statistics.pluginhelp
	case "debug":
		version = debug.pluginhelp
	case "link":
		version = link.pluginhelp
	case "games":
		version = games.pluginhelp
	}
	return
}

func init() {

}
