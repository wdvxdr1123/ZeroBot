package ZeroBot

// 插件信息
type PluginInfo struct {
	Author     string // 作者
	PluginName string // 插件名
	Version    string // 版本
	Details    string // 插件说明
}

var PluginPool []IPlugin

// 插件
type IPlugin interface {
	// 获取插件信息
	GetPluginInfo() PluginInfo
	// 开启工作
	Start()
}

func RegisterPlugin(plugin IPlugin) {
	PluginPool = append(PluginPool, plugin)
}
