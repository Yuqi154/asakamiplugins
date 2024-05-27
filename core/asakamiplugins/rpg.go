package asakamiplugins

import (
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type rpggameclientstruct struct {
	adress     string
	clientauth string
	expiretime time.Time
}

type rpgstruct struct {
	pluginversion string
	pluginhelp    string
	NewClient     func() rpggameclientstruct
}

type rpggamestruct struct {
	rpgcommon rpgcommonstruct
	rpgaction rpgactionstruct
	rpgbattle rpgbattlestruct
	rpgstatus rpgstatusstruct
	rpgloot   rpglootstruct
	rpgmap    rpgmapstruct
	rpgstore  rpgstorestruct
	rpgui     rpguistruct
	rpgserver rpgserverstruct
	rpgclient rpgclientstruct
}

type rpgcommonstruct struct {
	initrpg func(ctx *zero.Ctx) bool
}

type rpgactionstruct struct {
	// action
	moveto  func(x float64, y float64) bool
	explore func() bool
}

type rpgbattlestruct struct {
	// battle
	attack  func() bool
	defence func() bool
	skill   func() bool
	item    func() bool
	escape  func() bool
}

type rpgstatusstruct struct {
	// status
	getlevel        func() int
	getexp          func() float64
	getmaxhp        func() float64
	getmaxmp        func() float64
	gethp           func() float64
	getmp           func() float64
	getattack       func() float64
	getdefence      func() float64
	getmagicattack  func() float64
	getmagicdefence func() float64
	getspeed        func() float64
	getluck         func() float64
	getgold         func() int
	getskillpoint   func() int
	getstatuspoint  func() int
}

type rpglootstruct struct {
	// loot
	exploreloot func() bool
	battleloot  func() bool
}

type rpgmapstruct struct {
	// map
	initmap  func() bool
	loadmap  func(mapid int) bool
	syncmap  func() bool
	savemap  func() bool
	creature creaturestruct
}

type creaturestruct struct {
	build   func(buildingid int) bool
	upgrade func(buildingid int) bool
	destroy func(buildingid int) bool
}

type rpgstorestruct struct {
	// store
	buy       func(itemid int, quantity float64) bool
	sell      func(itemid int, quantity float64) bool
	list      func() bool
	search    func(itemid int) bool
	awaitsell func(itemid int, quantity float64, price float64) bool
	awaitbuy  func(itemid int, quantity float64, price float64) bool
}

type rpguistruct struct {
	// ui
	statusui   func() bool
	mapui      func() bool
	backpackui func() bool
	storeui    func() bool
	questui    func() bool
	upperui    func() bool
}

type rpgserverstruct struct {
	// server
	initserver func() bool
	listen     func(port int) bool
}

type rpgclientstruct struct {
	// client
	initclient func() bool
	connect    func(ip string, port int) bool
}

var rpgv = rpgstruct{
	pluginversion: "0.0.0a",
	pluginhelp:    "- 浅上RPG插件",
	NewClient: func() rpggameclientstruct {
		return rpggameclientstruct{
			adress:     "",
			clientauth: "",
			expiretime: time.Now().Add(time.Hour * 24 * 7),
		}
	},
}

var rpggamev = rpggamestruct{
	rpgcommon: rpgcommonv,
	rpgaction: rpgactionv,
	rpgbattle: rpgbattlev,
	rpgstatus: rpgstatusv,
	rpgloot:   rpglootv,
	rpgmap:    rpgmapv,
	rpgstore:  rpgstorev,
	rpgui:     rpguiv,
	rpgserver: rpgserverv,
	rpgclient: rpgclientv,
}

func Getrpg() rpgstruct {
	return rpgv
}

func initrpgfunc(ctx *zero.Ctx) bool {
	if ctx == nil {
		return false
	}
	uid := link.Getrealid(ctx)
	ctx.SendChain(message.Text("将初始化您的角色，是否继续？"))
	recv1, cancel1 := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^(是|否)$`), zero.CheckUser(uid)).Repeat()
	defer cancel1()
	for {
		select {
		case <-time.After(time.Second * 120):
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("等待超时,取消初始化")))
			return false
		case e := <-recv1:
			nextcmd := e.Event.Message.String()
			if nextcmd == "是" {
				goto next
			}
			if nextcmd == "否" {
				return false
			}
		}
	}
next:

	ctx.SendChain(message.Text("请选择您的性别：男/女"))
	recv2, cancel2 := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^(男|女)$`), zero.CheckUser(uid)).Repeat()
	defer cancel2()
	for {
		select {
		case <-time.After(time.Second * 120):
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("等待超时,取消初始化")))
			return false
		case e := <-recv2:
			nextcmd := e.Event.Message.String()
			if nextcmd == "男" {
				//
				goto next2
			}
			if nextcmd == "女" {
				//
				goto next2
			}
		}
	}
next2:

	ctx.SendChain(message.Text("请选择您的体型：矮小->高大(1-5)"))
	recv3, cancel3 := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^\d$`), zero.CheckUser(uid)).Repeat()
	defer cancel3()
	for {
		select {
		case <-time.After(time.Second * 120):
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("等待超时,取消初始化")))
			return false
		case e := <-recv3:
			nextcmd := e.Event.Message.String()
			switch nextcmd {
			case "1":
				//
				goto next3
			case "2":
				//
				goto next3
			case "3":
				//
				goto next3
			case "4":
				//
				goto next3
			case "5":
				//
				goto next3
			default:
				ctx.SendChain(message.Text("输入非法，请选择您的体型：矮小->高大(1-5)"))
			}
		}
	}
next3:

	return false
}

var rpgclientv = rpgclientstruct{
	initclient: func() bool {
		return true
	},
	connect: func(ip string, port int) bool {
		return true
	},
}

var rpgcommonv = rpgcommonstruct{
	initrpg: initrpgfunc,
}

var rpgactionv = rpgactionstruct{
	moveto: func(x float64, y float64) bool {
		return true
	},
	explore: func() bool {
		return true
	},
}

var rpgstatusv = rpgstatusstruct{
	getlevel: func() int {
		return 0
	},
	getexp: func() float64 {
		return 0
	},
	getmaxhp: func() float64 {
		return 0
	},
	getmaxmp: func() float64 {
		return 0
	},
	gethp: func() float64 {
		return 0
	},
	getmp: func() float64 {
		return 0
	},
	getattack: func() float64 {
		return 0
	},
	getdefence: func() float64 {
		return 0
	},
	getmagicattack: func() float64 {
		return 0
	},
	getmagicdefence: func() float64 {
		return 0
	},
	getspeed: func() float64 {
		return 0
	},
	getluck: func() float64 {
		return 0
	},
	getgold: func() int {
		return 0
	},
	getskillpoint: func() int {
		return 0
	},
	getstatuspoint: func() int {
		return 0
	},
}

var rpglootv = rpglootstruct{
	exploreloot: func() bool {
		return true
	},
	battleloot: func() bool {
		return true
	},
}

var rpgmapv = rpgmapstruct{
	initmap: func() bool {
		return true
	},
	loadmap: func(mapid int) bool {
		return true
	},
	syncmap: func() bool {
		return true
	},
	savemap: func() bool {
		return true
	},
	creature: creaturestruct{
		build: func(buildingid int) bool {
			return true
		},
		upgrade: func(buildingid int) bool {
			return true
		},
		destroy: func(buildingid int) bool {
			return true
		},
	},
}

var rpgstorev = rpgstorestruct{
	buy: func(itemid int, quantity float64) bool {
		return true
	},
	sell: func(itemid int, quantity float64) bool {
		return true
	},
	list: func() bool {
		return true
	},
	search: func(itemid int) bool {
		return true
	},
	awaitsell: func(itemid int, quantity float64, price float64) bool {
		return true
	},
	awaitbuy: func(itemid int, quantity float64, price float64) bool {
		return true
	},
}

var rpguiv = rpguistruct{
	statusui: func() bool {
		return true
	},
	mapui: func() bool {
		return true
	},
	backpackui: func() bool {
		return true
	},
	storeui: func() bool {
		return true
	},
	questui: func() bool {
		return true
	},
	upperui: func() bool {
		return true
	},
}

var rpgserverv = rpgserverstruct{
	initserver: func() bool {
		return true
	},
	listen: func(port int) bool {
		return true
	},
}

var rpgbattlev = rpgbattlestruct{
	attack: func() bool {
		return true
	},
	defence: func() bool {
		return true
	},
	skill: func() bool {
		return true
	},
	item: func() bool {
		return true
	},
	escape: func() bool {
		return true
	},
}

func (c rpggameclientstruct) connect() error {
	return nil
}

func (c rpggameclientstruct) disconnect() error {
	return nil
}

func (c rpggameclientstruct) login(uid int64, auth string) (gameplayerstruct, error) {
	return gameplayerstruct{uid: 0, expiretime: time.Now()}, nil
}

func (c rpggameclientstruct) logout(uid int64) error {
	return nil
}

func (c rpggameclientstruct) register(uid int64) string {
	return ""
}

func (c rpggameclientstruct) Command(str string) (gameresponsestruct, error) {
	rpggamev.rpgcommon.initrpg(nil)
	return gameresponsestruct{statuscode: 0, statusmsg: ""}, nil
}

type gameplayerstruct struct {
	uid        int64
	expiretime time.Time
}

type gameresponsestruct struct {
	statuscode int
	statusmsg  string
}

var client = rpggameclientstruct{
	adress:     "",
	clientauth: "",
	expiretime: time.Now().Add(time.Hour * 24 * 7),
}

func init() {
	client.connect()
	client.login(0, "")
	client.register(0)
	client.Command("")
	client.logout(0)
	client.disconnect()
}
