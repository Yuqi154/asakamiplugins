package asakamiplugins

type gamesstruct struct {
	pluginversion string
	pluginhelp    string
	foxfolder     foxfolderstruct
}

var gamesv = gamesstruct{
	pluginversion: "0.1.0",
	pluginhelp:    "- 浅上游戏插件",
	foxfolder:     Getfoxfolder(),
}

// Getgames 导出变量
func Getgames() gamesstruct {
	return gamesv
}
