package zero

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"testing"
)

type mockAPICaller struct{}

func (m mockAPICaller) CallApi(request APIRequest) (APIResponse, error) {
	return APIResponse{
		Status:  "",
		Data:    gjson.Result{},
		Msg:     "",
		Wording: "",
		RetCode: 0,
		Echo:    0,
	}, nil
}

// copy from extension.PatternModel
type PatternModel struct {
	Matched []*PatternParsed `zero:"pattern_matched"`
}

// Test Match
func TestText(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.MessageSegment{message.Text("haha")}, NewPattern().Text("haha"), true},
		{[]message.MessageSegment{message.Text("aaa")}, NewPattern().Text("not match"), false},
		{[]message.MessageSegment{message.Image("not a image")}, NewPattern().Text("not match"), false},
		{[]message.MessageSegment{message.At(114514)}, NewPattern().Text("not match"), false},
		{[]message.MessageSegment{message.Text("你说的对但是ZeroBot-Plugin 是 ZeroBot 的 实用插件合集")}, NewPattern().Text("实用插件合集"), true},
		{[]message.MessageSegment{message.Text("你说的对但是ZeroBot-Plugin 是 ZeroBot 的 实用插件合集")}, NewPattern().Text("nonono"), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := &Ctx{Event: &Event{Message: v.msg}}
			rule := PatternRule(v.pattern)
			out := rule(ctx)
			assert.Equal(t, v.expected, out)
		})
	}
}

func TestImage(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.MessageSegment{message.Text("haha")}, NewPattern().Image(), false},
		{[]message.MessageSegment{message.Text("haha"), message.Image("not a image")}, NewPattern().Image().Image(), false},
		{[]message.MessageSegment{message.Text("haha"), message.Image("not a image")}, NewPattern().Text("haha").Image(), true},
		{[]message.MessageSegment{message.Image("not a image")}, NewPattern().Image(), true},
		{[]message.MessageSegment{message.Image("not a image"), message.Image("not a image")}, NewPattern().Image(), false},
		{[]message.MessageSegment{message.Image("not a image"), message.Image("not a image")}, NewPattern().Image().Image(), true},
		{[]message.MessageSegment{message.Image("not a image"), message.Image("not a image")}, NewPattern().Image().Image().Image(), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := &Ctx{Event: &Event{Message: v.msg}}
			rule := PatternRule(v.pattern)
			out := rule(ctx)
			assert.Equal(t, v.expected, out)
		})
	}
}

func TestAt(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.MessageSegment{message.Text("haha")}, NewPattern().At(), false},
		{[]message.MessageSegment{message.Image("not a image")}, NewPattern().At(), false},
		{[]message.MessageSegment{message.At(114514)}, NewPattern().At(), true},
		{[]message.MessageSegment{message.At(114514)}, NewPattern().At("1919810"), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := &Ctx{Event: &Event{Message: v.msg}}
			rule := PatternRule(v.pattern)
			out := rule(ctx)
			assert.Equal(t, v.expected, out)
		})
	}
}

func TestReply(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.MessageSegment{message.Text("haha")}, NewPattern().Reply(), false},
		{[]message.MessageSegment{message.Image("not a image")}, NewPattern().Reply(), false},
		{[]message.MessageSegment{message.At(1919810), message.Reply(12345)}, NewPattern().Reply().At(), false},
		{[]message.MessageSegment{message.Reply(12345), message.At(1919810)}, NewPattern().Reply().At(), true},
		{[]message.MessageSegment{message.Reply(12345)}, NewPattern().Reply(), true},
		{[]message.MessageSegment{message.Reply(12345), message.At(1919810)}, NewPattern().Reply(), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := &Ctx{Event: &Event{Message: v.msg}, caller: APICaller(&mockAPICaller{})}
			rule := PatternRule(v.pattern)
			out := rule(ctx)
			assert.Equal(t, v.expected, out)
		})
	}
}
