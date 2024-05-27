package asakamiplugins

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/otiai10/openaigo"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	gptPluginVersion = "0.3.0"
	gptHelp          = "- 冬香酱写的gpt调用！"
	apiKey           = "api_key"
)

var filepath = engine.DataFolder() + "gpt/"

type gptstruct struct {
	pluginversion string
	pluginhelp    string
	run           func(ctx *zero.Ctx)
	reply         func(ctx *zero.Ctx)
}

var gptv = gptstruct{
	pluginversion: gptPluginVersion,
	pluginhelp:    gptHelp,
	run:           gptrun,
	reply:         gptreply,
}

// Getgpt 导出变量
func Getgpt() gptstruct {
	return gptv
}

func gptrun(ctx *zero.Ctx) {
	var gid = ctx.Event.GroupID

	prefix := regexp.MustCompile("^#GPT")
	text := ctx.Event.Message.ExtractPlainText()
	method := prefix.ReplaceAllString(text, "")
	switch method {
	case "查看历史记录":
		his, err := history("get", int(gid), []openaigo.Message{}, filepath+"history/")
		if err != nil {
			log.Print(err)
		}
		var res string
		for _, v := range his {
			res += fmt.Sprint(v.Role, ":", v.Content, "\n")
		}
		ctx.SendGroupMessage(gid, message.Text(res))
	case "清空历史记录":
		his, err := history("save", int(gid), []openaigo.Message{}, filepath+"history/")
		if err != nil {
			log.Print(err)
		}
		if len(his) != 0 {
			log.Print(errors.New("清空历史记录失败"))
		}
		ctx.SendGroupMessage(gid, message.Text("已清空历史记录"))
	case "查看预设":
		presets := getpresetlist(filepath)
		var res string = "预设列表:\n"
		for _, v := range presets {
			res += fmt.Sprint(" - ", v, "\n")
		}
		ctx.SendGroupMessage(gid, message.Text(res))
	default:
		if strings.Contains(method, "加载预设") {

			prefix := regexp.MustCompile("^加载预设")
			preset := prefix.ReplaceAllString(method, "")
			resa, err := fpreset(preset, filepath, gid)
			if err != nil {
				log.Print(err)
			}
			if resa == "" {
				ctx.SendGroupMessage(gid, message.Text("加载预设"+preset+"失败"))
			}
			ctx.SendGroupMessage(gid, message.Text(resa))
		} else {
			ctx.SendGroupMessage(gid, message.Text("未知指令"))
		}
	}
}

func gptreply(ctx *zero.Ctx) {

	var gid = ctx.Event.GroupID

	text := ctx.Event.Message.ExtractPlainText()

	//log.Print("[gptapi] loading")
	his, err := history("get", int(gid), []openaigo.Message{}, filepath+"history/")
	if err != nil {
		log.Print(err)
	}
	req := append(his, openaigo.Message{Role: "user", Content: text})
	res, resa := gpta(req)
	his, err = history("save", int(gid), res, filepath+"history/")
	if err != nil {
		log.Print(err)
	}
	if his == nil {
		log.Print(errors.New("保存历史记录失败"))
	}
	ctx.SendGroupMessage(gid, message.Text("GPT:"+resa))
}

func gpta(req []openaigo.Message) ([]openaigo.Message, string) {

	//proxyURL, _ := url.Parse("http://127.0.0.1:7890")
	client := openaigo.NewClient(apiKey)
	//client.HTTPClient.Transport = &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	Proxy:           http.ProxyURL(proxyURL),
	//}
	request := openaigo.ChatRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    req,
		N:           1,
		Temperature: 0.7,
		TopP:        0.8,
		Stream:      false,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel() // 释放资源
	//log.Print(request)
	response, err := client.Chat(ctx, request)
	if err != nil {
		log.Print(err)
	}
	//log.Print(response)
	res := append(req, response.Choices[0].Message)
	return res, response.Choices[0].Message.Content
}

func history(method string, id int, his []openaigo.Message, filepath string) ([]openaigo.Message, error) {
	switch method {
	case "get":
		//从id.his读取id的历史记录
		file, err := os.Open(filepath + fmt.Sprint(id, ".his"))
		if err != nil {
			return []openaigo.Message{}, err
		}
		defer file.Close()
		//读取json文件
		reader := bufio.NewReader(file)
		//读取文件内容
		content, err := ioutil.ReadAll(reader)
		if err != nil {
			return []openaigo.Message{}, err
		}
		//解析json
		var history []openaigo.Message
		err = json.Unmarshal(content, &history)
		if err != nil {
			return []openaigo.Message{}, err
		}
		return history, nil
	case "save":
		//将id的历史记录保存到id.his
		file, err := os.OpenFile(filepath+fmt.Sprint(id, ".his"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return []openaigo.Message{}, err
		}
		defer file.Close()
		//将历史记录转换为json
		content, err := json.Marshal(his)
		if err != nil {
			return []openaigo.Message{}, err
		}
		//写入文件
		_, err = file.Write(content)
		if err != nil {
			return []openaigo.Message{}, err
		}
		return []openaigo.Message{}, nil

	default:
		return []openaigo.Message{}, errors.New("unknown method")
	}
}

func fpreset(presetf string, filepath string, gid int64) (string, error) {
	//从preset.prs读取preset
	file, err := os.Open(filepath + "preset/" + fmt.Sprint(presetf, ".preset"))
	if err != nil {
		return "", err
	}
	defer file.Close()
	log.Print("打开文件成功")
	//读取json文件
	reader := bufio.NewReader(file)
	//读取文件内容
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	//解析json
	var preset []openaigo.Message
	log.Print("解析json")
	err = json.Unmarshal(content, &preset)
	if err != nil {
		return "", err
	}
	//调用gpt
	log.Print("调用gpt")
	res, cont := gpta(preset)
	his, err := history("save", int(gid), res, filepath+"history/")
	//保存历史记录
	log.Print("保存历史记录")
	if err != nil {
		log.Print(err)
	}
	if his == nil {
		log.Print(errors.New("保存历史记录失败"))
	}
	return cont, nil
}

func getpresetlist(filepath string) []string {
	//读取preset目录
	files, err := ioutil.ReadDir(filepath + "preset/")
	if err != nil {
		log.Print(err)
	}
	//遍历preset目录
	var presetlist []string
	for _, file := range files {
		//获取文件名
		filename := file.Name()
		//获取文件后缀
		filetype := path.Ext(filename)
		//判断文件后缀是否为.preset
		if filetype == ".preset" {
			//获取文件名（不包含后缀）
			filenameOnly := strings.TrimSuffix(filename, filetype)
			//将文件名加入presetlist
			presetlist = append(presetlist, filenameOnly)
		}
	}
	return presetlist
}
