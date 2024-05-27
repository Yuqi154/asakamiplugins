# AsakamiPlugins

这是由asakamifuyuka制作的插件，包括大量功能。

## 插件列表

### core/asakamiplugins

<details>
  <summary> `core.go` 组件内部注册，帮助信息的显示</summary>

  - [x] #浅上/#ASAKAMI help/status/version/list/[插件名] help/[插件名] status/[插件名] version

</details>

<details>
  <summary> `entry.go` 插件入口，负责注册engine</summary>

  - 内部调用，无触发
</details>
<details>
  <summary> `chat.go` `smmodel.go` 聊天功能，包括睡眠管理，骰娘</summary>

  与zerobotplugins的`sleep_manage.go`功能重复

  - [x] 早安/晚安(各种语言和变体)

  - [x] roll [tag] [times]d[sides]\(!)/[need]  
    tag: 标签 (必须) times: 次数 (必须) sides: 面数 (4/6/8/10/12/20/100) need: 需求 (必须) (!): 忽略面数检查
</details>
<details>
  <summary> `aipaint.go` 与stable-diffusion的AI画图交互</summary>

  - [x] #春生画图+prompt:{...}+negative_prompt:{...}+steps=..+h=...+w=...+cfg_scale=...+seed=...  
    默认api地址[https://127.0.0.1:7860](https://127.0.0.1:7860)
</details>
<details>
  <summary> `backpack.go` 背包</summary>

  - 内部调用，无触发
    AsakamiPlugins的核心数据库
</details>
<details>
  <summary> `debug.go` 调试功能</summary>

  - [x] #debug#uid  
    显示uid(用户当前id)gid(群组id)realuid(用户真实id)

  - [x] #debug#runtime  
    显示运行时间
</details>
<details>
  <summary> `economy.go` 经济系统</summary>

  - 内部调用，无触发
</details>
<details>
  <summary> `foxfolder.go` 叠狐狸</summary>

  - [x] #叠狐狸 加入/分析/列表

</details>
<details>
  <summary> `gameentry.go` 游戏入口</summary>

  - 内部调用，无触发
</details>
<details>
  <summary> `gpt.go` 与gpt的交互</summary>

  - [x] #GPT / @机器人  
    默认不启用
</details>
<details>
  <summary> `heartbeat.go` 心跳</summary>

  - 内部调用，无触发  
    与runtime分析相关
</details>
<details>
  <summary> `link.go` 绑定主账户</summary>

  - [x] #绑定QQ / #绑定KOOK  
    需要确认并使用绑定目标账户发送验证码
</details>
<details>
  <summary> `pay.go` 稻荷宝</summary>

  - [x] #稻荷宝 [help/帮助]
  - [x] #稻荷宝 [pay/支付] [u[uid]/@] 金额
  - [x] #稻荷宝 [balance/余额]
  - [x] #稻荷宝 [金额] [jpy/dhc/dyc] [exchange/兑换] [jpy/dhc/dyc] 
  - [x] #稻荷宝 [rate/汇率] 
</details>
<details>
  <summary> `rcon.go` RCON</summary>

  - [x] #rcon send#[command]

    发送RCON命令,默认服务器[127.0.0.1:25575]()  
    需修改允许的QQuid
</details>
<details>
  <summary>(未完成) `rpg.go` `rpgui.go` RPG</summary>
</details>
<details>
  <summary> `sign_in.go` `signinmodel.go` `signinui.go` 签到</summary>
与zerobotplugins的`score.go`功能重复

  - [x] 签到
  - [x] 获得(取)签到背景[@xxx] | 获得(取)签到背景
  - [x] 设置签到预设(0~3)
  - [x] 查看等级排名
  - [x] 签到数据
</details>
<details>
  <summary> `statistic.go` 统计</summary>

  - 内部调用，无触发  
</details>
<details>
  <summary>(未完成) `wordcloud.go` 词云</summary>
</details>

### modify

修改的部分zerobotplugins以兼容本插件数据库