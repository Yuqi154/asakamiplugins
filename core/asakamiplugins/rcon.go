package asakamiplugins

import (
	"regexp"
	"strings"

	mcrcon "github.com/Kelwing/mc-rcon"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type Rcon struct {
	Addr string
	Pass string
}

func (r *Rcon) Send(cmd string) (resp string, err error) {
	conn := new(mcrcon.MCConn)
	err = conn.Open(r.Addr, r.Pass)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	err = conn.Authenticate()
	if err != nil {
		return "", err
	}

	resp, err = conn.SendCommand(cmd)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (r *Rcon) run(ctx *zero.Ctx) {
	text := ctx.Event.Message.ExtractPlainText()
	if text == "rcon" {
		ctx.SendGroupMessage(ctx.Event.GroupID, message.Text("请输入rcon命令"))
		return
	}
	if !(ctx.Event.UserID == 1123489751) {
		ctx.SendGroupMessage(ctx.Event.GroupID, message.Text("你没有权限使用这个命令"))
		return
	}
	rcon := new(Rcon)
	rcon.Addr = "localhost:25575"
	rcon.Pass = "fuyuka"
	var sendstr string
	if strings.Contains(text, "send#") || strings.Contains(text, "SEND#") {
		sendrep := regexp.MustCompile(`(send|SEND)#(.*)`)
		sendstr = sendrep.FindStringSubmatch(text)[0]
		sendstr = strings.Replace(sendstr, "send#", "", 1)
	}
	if sendstr == "" {
		ctx.SendGroupMessage(ctx.Event.GroupID, message.Text("请输入rcon命令"))
		return
	}
	resp, err := rcon.Send(sendstr)
	if err != nil {
		ctx.SendGroupMessage(ctx.Event.GroupID, message.Text("发送失败"))
		return
	}
	ctx.SendGroupMessage(ctx.Event.GroupID, message.Text(resp))
}

func GetRcon() Rcon {
	return Rcon{}
}
