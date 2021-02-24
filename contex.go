package zero

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/modern-go/reflect2"
)

// Ctx represents the Context which hold the event.
// 代表上下文
type Ctx struct {
	ma    *Matcher
	Event *Event
	State State
}

// GetMatcher ...
func (ctx *Ctx) GetMatcher() *Matcher {
	return ctx.ma
}

// decoder 反射获取的数据
type decoder []struct {
	offset uintptr
	t      reflect2.Type
	key    string
}

// decoder 缓存
var decoderCache = sync.Map{}

// Parse 将 Ctx.State 映射到结构体
func (ctx *Ctx) Parse(model interface{}) (err error) {
	var (
		ty2      = reflect2.TypeOf(model)
		modelDec decoder
	)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse state error: %v", r)
		}
	}()
	dec, ok := decoderCache.Load(ty2)
	if ok {
		modelDec = dec.(decoder)
	} else {
		t := ty2.(reflect2.PtrType).Elem().(reflect2.StructType)
		modelDec = decoder{}
		for i := 0; i < t.NumField(); i++ {
			t1 := t.Field(i)
			if key, ok := t1.Tag().Lookup("zero"); ok {
				modelDec = append(modelDec, struct {
					offset uintptr
					t      reflect2.Type
					key    string
				}{
					t:      t1.Type(),
					offset: t1.Offset(),
					key:    key,
				})
			}
		}
		decoderCache.Store(ty2, modelDec)
	}
	for i := range modelDec { // decoder类型非小内存，无法被编译器优化为快速拷贝
		modelDec[i].t.UnsafeSet(
			unsafe.Pointer(uintptr(reflect2.PtrOf(model))+modelDec[i].offset),
			reflect2.PtrOf(ctx.State[modelDec[i].key]),
		)
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
func (ctx *Ctx) Send(message interface{}) int64 {
	if ctx.Event.GroupID != 0 {
		return SendGroupMessage(ctx.Event.GroupID, message)
	}
	return SendPrivateMessage(ctx.Event.UserID, message)
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
	return (<-ctx.FutureEvent("message", ctx.CheckSession()).Next()).RawMessage
}
