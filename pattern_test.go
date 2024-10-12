package zero

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"testing"
)

type mockAPICaller struct{}

func (m mockAPICaller) CallApi(_ APIRequest) (APIResponse, error) {
	return APIResponse{
		Status:  "",
		Data:    gjson.Result{},
		Msg:     "",
		Wording: "",
		RetCode: 0,
		Echo:    0,
	}, nil
}
func fakeCtx(msg message.Message) *Ctx {
	ctx := &Ctx{Event: &Event{Message: msg}, State: map[string]interface{}{}, caller: mockAPICaller{}}
	return ctx
}

// copy from extension.PatternModel
type PatternModel struct {
	Matched []*PatternParsed `zero:"pattern_matched"`
}

// Test Match
func TestPattern_Text(t *testing.T) {
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
			ctx := fakeCtx(v.msg)
			rule := PatternRule(v.pattern)
			out := rule(ctx)
			assert.Equal(t, out, v.expected)
		})
	}
}

func TestPattern_Image(t *testing.T) {
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
			ctx := fakeCtx(v.msg)
			rule := PatternRule(v.pattern)
			out := rule(ctx)
			assert.Equal(t, out, v.expected)
		})
	}
}

func TestPattern_At(t *testing.T) {
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
			ctx := fakeCtx(v.msg)
			rule := PatternRule(v.pattern)
			out := rule(ctx)
			assert.Equal(t, out, v.expected)
		})
	}
}

func TestPattern_Reply(t *testing.T) {
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
			ctx := fakeCtx(v.msg)
			rule := PatternRule(v.pattern)
			out := rule(ctx)
			assert.Equal(t, out, v.expected)
		})
	}
}
func TestPatternParsed_Gets(t *testing.T) {
	assert.Equal(t, []string{"gaga"}, PatternParsed{Valid: true, Value: []string{"gaga"}}.GetText())
	assert.Equal(t, "image", PatternParsed{Valid: true, Value: "image"}.GetImage())
	assert.Equal(t, "reply", PatternParsed{Valid: true, Value: "reply"}.GetReply())
	assert.Equal(t, "114514", PatternParsed{Valid: true, Value: "114514"}.GetAt())
}
func TestPattern_SetOptional(t *testing.T) {
	assert.Panics(t, func() {
		NewPattern().SetOptional()
	})
	tests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected []PatternParsed
	}{
		{[]message.MessageSegment{message.Text("/do it")}, NewPattern().Text("/(do) (.*)").At().SetOptional(true), []PatternParsed{
			{
				Valid: true,
			}, {
				Valid: false,
			},
		}},
		{[]message.MessageSegment{message.Text("/do it")}, NewPattern().Text("/(do) (.*)").At().SetOptional(false), []PatternParsed{}},
		{[]message.MessageSegment{message.Text("happy bear"), message.At(114514)}, NewPattern().Reply().SetOptional().Text(".+").SetOptional().At().SetOptional(false), []PatternParsed{
			{
				Valid: false,
			},
			{
				Valid: true,
			},
			{
				Valid: true,
			},
		}},
		{[]message.MessageSegment{message.Text("happy bear"), message.At(114514)}, NewPattern().Image().SetOptional().Image().SetOptional().Image().SetOptional(), []PatternParsed{ // why you do this
			{
				Valid: false,
			},
			{
				Valid: false,
			},
			{
				Valid: false,
			},
		}},
	}
	for i, v := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := PatternRule(v.pattern)
			matched := rule(ctx)
			if !matched {
				assert.Equal(t, 0, len(v.expected))
				return
			}
			parsed := &PatternModel{}
			err := ctx.Parse(parsed)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, len(v.expected), len(parsed.Matched))
			for i := range parsed.Matched {
				assert.Equal(t, v.expected[i].Valid, parsed.Matched[i].Valid)
			}
		})
	}
}

// Test Parse
func TestAllParse(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected []PatternParsed
	}{
		{[]message.MessageSegment{message.Text("test haha test"), message.At(123)}, NewPattern().Text("((ha)+)").At(), []PatternParsed{
			{
				Valid: true,
				Value: []string{"haha", "haha", "ha"},
			}, {
				Valid: true,
				Value: int64(123),
			},
		}},
		{[]message.MessageSegment{message.Text("haha")}, NewPattern().Text("(h)(a)(h)(a)"), []PatternParsed{
			{
				Valid: true,
				Value: []string{"haha", "h", "a", "h", "a"},
			},
		}},
		{[]message.MessageSegment{message.Reply("fake reply"), message.Image("fake image"), message.At(999), message.At(124), message.Text("haha")}, NewPattern().Reply().Image().At().At("124").Text("(h)(a)(h)(a)"), []PatternParsed{

			{
				Valid: true,
				Value: "fake reply",
			},
			{
				Valid: true,
				Value: "fake image",
			},
			{
				Valid: true,
				Value: int64(999),
			},
			{
				Valid: true,
				Value: int64(124),
			},
			{
				Valid: true,
				Value: []string{"haha", "h", "a", "h", "a"},
			},
		}},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := PatternRule(v.pattern)
			matched := rule(ctx)
			parsed := &PatternModel{}
			err := ctx.Parse(parsed)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, true, matched)
			for i := range parsed.Matched {
				assert.Equal(t, v.expected[i].Valid, parsed.Matched[i].Valid)
				assert.Equal(t, v.expected[i].Value, parsed.Matched[i].Value)
				assert.Equal(t, &(v.msg[i]), parsed.Matched[i].Msg)
			}
		})
	}
}
