---
title: "First"
date: 2020-12-26T12:57:28+08:00
draft: false
---

## 开始使用

## 最小化实例

创建一个新项目，启动Bot只需要main中使用下面代码

```golang
func main() {
    zero.Run(zero.Option{
        Host:          "127.0.0.1", // cqhttp的ip地址
        Port:          "6700", // cqhttp的端口
        AccessToken:   "",
        NickName:      []string{"机器人的昵称"},
        CommandPrefix: "/", // 指令前缀
        SuperUsers:    []string{"123456"}, // 超级用户账号 一般填你自己的QQ号
    })
    select {} // 阻塞主goroutine, 防止退出程序
    // 如果你的机器人有使用数据库或者其他资源文件,可以使用下面方法阻塞
    /*
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    <-c
    Release() // 释放资源
    */
}
```

## 设置日志输出

在 ZeroBot 中使用了`sirupsen/logrus`来管理日志，但是并没有提供日志的模板，你可以自己定义日志输出模板，
如果想偷懒，这里给出一个简单的模板

```golang
import (
    log "github.com/sirupsen/logrus"
    easy "github.com/t-tomalak/logrus-easy-formatter"
)

func init() {
    log.SetFormatter(&easy.Formatter{
        TimestampFormat: "2006-01-02 15:04:05",
        LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
    })
    log.SetLevel(log.DebugLevel)
}
```

## 添加一个插件

得益于 GO111MODULE 和 `init` 函数，在ZeroBot里添加一个插件十分简单，如果想要添加`example`代码中的插件，只需要

```golang
import _ "github.com/wdvxdr1123/ZeroBot/example/music"
```

添加这个导入后即可为你的机器人添加点歌插件

## 编写一个插件

在 ZeroBot 中，插件是以接口的形式实现的

```golang
// IPlugin is the plugin of the ZeroBot
type IPlugin interface {
    // 获取插件信息
    GetPluginInfo() PluginInfo
    // 开启工作
    Start()
}
```

编写一个插件你需要新建一个文件夹，然后自定义一个类型，为它实现`IPlugin`接口,
并在`init`函数中，注册该插件,例如

```golang
package test

import (
    "github.com/wdvxdr1123/ZeroBot"
)

func init() {
    zero.RegisterPlugin(&testPlugin{}) // 注册插件
}

type testPlugin struct{}

func (_ *testPlugin) GetPluginInfo() zero.PluginInfo { // 返回插件信息
    return zero.PluginInfo{       // 插件信息本身没什么用，但是可以方便别人了解你写的插件
        Author:     "wdvxdr1123", // 作者
        PluginName: "test",       // 插件名 
        Version:    "0.1.0",      // 版本
        Details:    "这是一个测试插件", // 插件信息
    }
}

func (_ *testPlugin) Start() { // 插件主体
    panic("impl me")
}
```

### 编写一个复读插件

```golang
func (_ *testPlugin) Start() { // 插件主体
    zero.OnCommand("echo").Handle(handleEcho) // 注册一个叫echo的指令，逻辑处理函数为 handleEcho
}

func handleEcho(_ *zero.Matcher, event zero.Event, state zero.State) zero.Response {
    zero.Send(event, state["args"]) // 发送echo的参数
    return zero.FinishResponse // 所有处理已完毕，返回Finish
}
```
