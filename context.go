package zero

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/wdvxdr1123/ZeroBot/message"
)

// Ctx represents the Context which hold the event.
// 代表上下文
type Ctx struct {
	ma     *Matcher
	Event  *Event
	State  State
	caller APICaller

	// lazy message
	once    sync.Once
	message string
}

// GetMatcher ...
func (ctx *Ctx) GetMatcher() *Matcher {
	return ctx.ma
}

// decoder 反射获取的数据
type decoder []dec

type dec struct {
	index int
	key   string
}

// decoder 缓存
var decoderCache = sync.Map{}

// Parse 将 Ctx.State 映射到结构体
func (ctx *Ctx) Parse(model interface{}) (err error) {
	var (
		rv       = reflect.ValueOf(model).Elem()
		t        = rv.Type()
		modelDec decoder
	)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse state error: %v", r)
		}
	}()
	d, ok := decoderCache.Load(t)
	if ok {
		modelDec = d.(decoder)
	} else {
		modelDec = decoder{}
		for i := 0; i < t.NumField(); i++ {
			t1 := t.Field(i)
			if key, ok := t1.Tag.Lookup("zero"); ok {
				modelDec = append(modelDec, dec{
					index: i,
					key:   key,
				})
			}
		}
		decoderCache.Store(t, modelDec)
	}
	for _, d := range modelDec { // decoder类型非小内存，无法被编译器优化为快速拷贝
		rv.Field(d.index).Set(reflect.ValueOf(ctx.State[d.key]))
	}
	return nil
}

// CheckSession 判断会话连续性
func (ctx *Ctx) CheckSession() Rule {
	return func(ctx2 *Ctx) bool {
		return ctx.Event.UserID == ctx2.Event.UserID &&
			ctx.Event.GroupID == ctx2.Event.GroupID // 私聊时GroupID为0，也相等
	}
}

// Send 快捷发送消息
func (ctx *Ctx) Send(msg interface{}) message.MessageID {
	event := ctx.Event
	if event.DetailType == "guild" {
		return message.NewMessageIDFromString(ctx.SendGuildChannelMessage(event.GuildID, event.ChannelID, msg))
	}
	if event.GroupID != 0 {
		return message.NewMessageIDFromInteger(ctx.SendGroupMessage(event.GroupID, msg))
	}
	return message.NewMessageIDFromInteger(ctx.SendPrivateMessage(event.UserID, msg))
}

// SendChain 快捷发送消息-消息链
func (ctx *Ctx) SendChain(msg ...message.MessageSegment) message.MessageID {
	return ctx.Send((message.Message)(msg))
}

// Echo 向自身分发虚拟事件
func (ctx *Ctx) Echo(response []byte) {
	processEvent(response, ctx.caller)
}

// FutureEvent ...
func (ctx *Ctx) FutureEvent(Type string, rule ...Rule) *FutureEvent {
	return ctx.ma.FutureEvent(Type, rule...)
}

// Get ..
func (ctx *Ctx) Get(prompt string) string {
	if prompt != "" {
		ctx.Send(prompt)
	}
	return (<-ctx.FutureEvent("message", ctx.CheckSession()).Next()).Event.RawMessage
}

// ExtractPlainText 提取消息中的纯文本
func (ctx *Ctx) ExtractPlainText() string {
	if ctx == nil || ctx.Event == nil || ctx.Event.Message == nil {
		return ""
	}
	return ctx.Event.Message.ExtractPlainText()
}

// Block 阻止后续触发
func (ctx *Ctx) Block() {
	ctx.ma.SetBlock(true)
}

// MessageString 字符串消息便于Regex
func (ctx *Ctx) MessageString() string {
	ctx.once.Do(func() {
		if ctx.Event != nil && ctx.Event.Message != nil {
			ctx.message = ctx.Event.Message.String()
		}
	})
	return ctx.message
}
