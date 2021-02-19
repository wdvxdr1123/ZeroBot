package zero

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
	// GetPluginInfo 获取插件信息
	GetPluginInfo() PluginInfo
	// Start 开启工作
	Start()
}

// RegisterPlugin register the plugin to ZeroBot
func RegisterPlugin(plugin IPlugin) {
	pluginPool = append(pluginPool, plugin)
}
