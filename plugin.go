package zero

import log "github.com/sirupsen/logrus"

// PluginInfo is the plugin's information
type PluginInfo struct {
	Author     string // 作者
	PluginName string // 插件名
	Version    string // 版本
	Details    string // 插件说明
}

var pluginPool []IPlugin

// IPlugin is the plugin of the ZeroBot
type IPlugin interface {
	// 获取插件信息
	GetPluginInfo() PluginInfo
	// 开启工作
	Start()
}

// RegisterPlugin register the plugin to ZeroBot
func RegisterPlugin(plugin IPlugin) {
	info := plugin.GetPluginInfo()
	log.Infof(
		"加载插件: %v [作者] %v [版本] %v [说明] %v",
		info.PluginName,
		info.Author,
		info.Version,
		info.Details,
	)
	pluginPool = append(pluginPool, plugin)
}
