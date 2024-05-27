package asakamiplugins

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/FloatTech/AnimeAPI/wallet"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 定义aipaint结构体
type aipaintstruct struct {
	pluginversion string
	pluginhelp    string
	run           func(ctx *zero.Ctx)
}

// 定义常量
const (
	aipaintpluginVersion = "0.3.1"
	aipainthelp          = "- 冬香酱写的aipaint \n用法:\n#春生画图+prompt:{...}+negative_prompt:{...}+steps=..+h=...+w=...+cfg_scale=...+seed=..."
)

// 定义aipaint结构体
var aipaintv = aipaintstruct{
	pluginversion: aipaintpluginVersion,
	pluginhelp:    aipainthelp,
	run:           aipaintrun,
}

// 导出aipaint结构体
func Getaipaint() aipaintstruct {
	return aipaintv
}

var (
	jsontmp = map[string]interface{}{
		"enable_hr":                            false,
		"denoising_strength":                   0,
		"firstphase_width":                     0,
		"firstphase_height":                    0,
		"hr_scale":                             2,
		"hr_upscaler":                          "string",
		"hr_second_pass_steps":                 0,
		"hr_resize_x":                          0,
		"hr_resize_y":                          0,
		"prompt":                               "",
		"seed":                                 -1,
		"subseed":                              -1,
		"subseed_strength":                     0,
		"seed_resize_from_h":                   -1,
		"seed_resize_from_w":                   -1,
		"sampler_name":                         "Euler a",
		"batch_size":                           1,
		"n_iter":                               1,
		"steps":                                20,
		"cfg_scale":                            7,
		"width":                                512,
		"height":                               512,
		"restore_faces":                        false,
		"tiling":                               false,
		"negative_prompt":                      "",
		"eta":                                  0,
		"s_churn":                              0,
		"s_tmax":                               0,
		"s_tmin":                               0,
		"s_noise":                              1,
		"override_settings_restore_afterwards": true,
		"sampler_index":                        "Euler",
		"script_name":                          "",
		"send_images":                          true,
		"save_images":                          false,
	}
)

func aipaintrun(ctx *zero.Ctx) {

	//gid := ctx.Event.GroupID
	uid := link.Getrealid(ctx)
	money := wallet.GetWalletOf(uid)

	text := ctx.Event.Message.ExtractPlainText()
	prefix := regexp.MustCompile("^#春生画图")
	r_text := prefix.ReplaceAllString(text, "")

	if strings.Contains(r_text, "help") || strings.Contains(r_text, "帮助") {
		ctx.SendChain(message.Text(aipainthelp))
		return
	}
	//帮助

	//测试http://localhost:7861链接
	_, err := http.Get("https://127.0.0.1:7860")
	if err != nil {
		log.Print(err)
		//ctx.SendChain(message.Text("链接错误,可能是stable diffusion未启动"))//
		//return
	}
	//检查链接

	prompt := getindex(r_text, "prompt:{", "}")
	negative_prompt := getindex(r_text, "negative_prompt:{", "}")
	stepss := getindex2(r_text, "steps=", "\\d*")
	hs := getindex2(r_text, "(h|H)=", "\\d*")
	ws := getindex2(r_text, "(w|W)=", "\\d*")
	seeds := getindex2(r_text, "seed=", "[\\d-]*")
	cfg_scales := getindex2(r_text, "(cfg|CFG)_scale=", "[\\d-]*")
	//获取参数

	steps, err := strconv.ParseInt(stepss, 0, 8)
	h, err2 := strconv.ParseInt(hs, 0, 16)
	w, err3 := strconv.ParseInt(ws, 0, 16)
	seed, err4 := strconv.ParseInt(seeds, 0, 64)
	cfg_scale, err5 := strconv.ParseFloat(cfg_scales, 64)
	//转换参数

	if !strings.Contains(r_text, "steps") {
		steps = 20
	}
	if !(strings.Contains(r_text, "h=") || strings.Contains(r_text, "H=")) {
		h = 512
	}
	if !(strings.Contains(r_text, "w=") || strings.Contains(r_text, "W=")) {
		w = 512
	}
	if !strings.Contains(r_text, "seed") {
		seed = -1
	}
	if !(strings.Contains(r_text, "cfg_scale") || strings.Contains(r_text, "CFG_scale")) {
		cfg_scale = 7
	}
	//检查参数

	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		log.Print(err)
		//ctx.SendChain( message.Text("参数错误"))
		//return
	}
	if steps < 1 || steps > 100 {
		ctx.SendChain(message.Text("步数越界，请检查steps="))
		return
	}
	if h%64 != 0 {
		//log.Print(h)
		ctx.SendChain(message.Text("图片高度错误，请检查h="))
		return
	}
	if h < 128 || h > 1024 {
		ctx.SendChain(message.Text("图片高度越界，请检查h="))
		return
	}
	if w%64 != 0 {
		//log.Print(w)
		ctx.SendChain(message.Text("图片宽度错误，请检查w="))
		return
	}
	if w < 128 || w > 1024 {
		ctx.SendChain(message.Text("图片宽度越界，请检查w="))
		return
	}
	if int(cfg_scale*10)%5 != 0 {
		//log.Print(w)
		ctx.SendChain(message.Text("cfg_scale错误，请检查cfg_scale="))
		return
	}
	//检查参数

	cost := float64(h / 256 * w / 256 * steps / 20)
	if money < cost {
		ctx.SendChain(message.Text("余额不足(但是可以继续使用)"))
	} else {
		wallet.InsertWalletOf(uid, -float64(cost))
		ctx.SendChain(message.Text("消费" + fmt.Sprint(cost) + "稻荷币，余额" + fmt.Sprint(money-cost) + "稻荷币"))
	}
	//扣钱

	ctx.SendChain(message.Text("少女祈祷中..."))

	jsontmp["prompt"] = prompt
	jsontmp["negative_prompt"] = negative_prompt
	jsontmp["steps"] = steps
	jsontmp["height"] = h
	jsontmp["width"] = w
	jsontmp["seed"] = seed
	jsontmp["cfg_scale"] = cfg_scale
	//设置参数

	jsonbyte, err := json.Marshal(jsontmp)
	if err != nil {
		ctx.SendChain(message.Text("json错误"))
		return
	}
	//转换json

	jsonpost := bytes.NewBuffer(jsonbyte)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err2 := client.Post("https://localhost:7860/sdapi/v1/txt2img", "application/json", jsonpost)
	//发送post请求

	if err2 != nil {
		ctx.SendChain(message.Text("post错误:" + err2.Error()))
		return
	}
	//检查post请求
	defer resp.Body.Close()
	//关闭请求
	body, _ := ioutil.ReadAll(resp.Body)
	//读取返回值
	var data map[string]interface{}
	//定义map
	json.Unmarshal(body, &data)
	//解析json
	images := data["images"]
	//获取图片
	str := fmt.Sprintf("%s", images)
	//转换图片
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "[", "")
	str = strings.ReplaceAll(str, "]", "")
	str = strings.ReplaceAll(str, "\\t", "")
	str = strings.ReplaceAll(str, "\\n", "")
	//处理图片
	img, err := base64.StdEncoding.DecodeString(str)
	//解码图片
	if err != nil {
		ctx.SendChain(message.Text("图片解析错误"))
		log.Printf("Error decoding2 string:  %s ", err.Error())
	}
	//检查图片
	ctx.SendChain(message.ImageBytes(img))
	//发送图片
}

// 获取字符串中间的内容
func getindex(text string, head string, tail string) string {

	prefix := regexp.MustCompile(head + ".*" + tail)
	prefix2 := regexp.MustCompile(head + "|" + tail)
	index0 := prefix.FindString(text)
	index := prefix2.ReplaceAllString(index0, "")
	return index

}

// 获取字符串中间的内容
func getindex2(text string, head string, tail string) string {

	prefix := regexp.MustCompile(head + tail)
	prefix2 := regexp.MustCompile(head)
	index0 := prefix.FindString(text)
	index := prefix2.ReplaceAllString(index0, "")
	return index
}
