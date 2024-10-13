package zero

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"testing"
)

type mockAPICaller struct{}

func (m mockAPICaller) CallAPI(_ APIRequest) (APIResponse, error) {
	return APIResponse{
		Status:  "",
		Data:    gjson.Parse(`{"message_id":"12345","sender":{"user_id":12345}}`), // just for reply cleaner
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
		{[]message.Segment{message.Text("haha")}, NewPattern().Text("haha"), true},
		{[]message.Segment{message.Text("aaa")}, NewPattern().Text("not match"), false},
		{[]message.Segment{message.Image("not a image")}, NewPattern().Text("not match"), false},
		{[]message.Segment{message.At(114514)}, NewPattern().Text("not match"), false},
		{[]message.Segment{message.Text("你说的对但是ZeroBot-Plugin 是 ZeroBot 的 实用插件合集")}, NewPattern().Text("实用插件合集"), true},
		{[]message.Segment{message.Text("你说的对但是ZeroBot-Plugin 是 ZeroBot 的 实用插件合集")}, NewPattern().Text("nonono"), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
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
		{[]message.Segment{message.Text("haha")}, NewPattern().Image(), false},
		{[]message.Segment{message.Text("haha"), message.Image("not a image")}, NewPattern().Image().Image(), false},
		{[]message.Segment{message.Text("haha"), message.Image("not a image")}, NewPattern().Text("haha").Image(), true},
		{[]message.Segment{message.Image("not a image")}, NewPattern().Image(), true},
		{[]message.Segment{message.Image("not a image"), message.Image("not a image")}, NewPattern().Image(), false},
		{[]message.Segment{message.Image("not a image"), message.Image("not a image")}, NewPattern().Image().Image(), true},
		{[]message.Segment{message.Image("not a image"), message.Image("not a image")}, NewPattern().Image().Image().Image(), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
			out := rule(ctx)
			assert.Equal(t, v.expected, out)
		})
	}
}

func TestPattern_At(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.Segment{message.Text("haha")}, NewPattern().At(), false},
		{[]message.Segment{message.Image("not a image")}, NewPattern().At(), false},
		{[]message.Segment{message.At(114514)}, NewPattern().At(), true},
		{[]message.Segment{message.At(114514)}, NewPattern().At(message.NewMessageIDFromString("1919810")), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
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
		{[]message.Segment{message.Text("haha")}, NewPattern().Reply(), false},
		{[]message.Segment{message.Image("not a image")}, NewPattern().Reply(), false},
		{[]message.Segment{message.At(1919810), message.Reply(12345)}, NewPattern().Reply().At(), false},
		{[]message.Segment{message.Reply(12345), message.At(1919810)}, NewPattern().Reply().At(), true},
		{[]message.Segment{message.Reply(12345)}, NewPattern().Reply(), true},
		{[]message.Segment{message.Reply(12345), message.At(1919810)}, NewPattern().Reply(), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
			out := rule(ctx)
			assert.Equal(t, out, v.expected)
		})
	}
}
func TestPattern_ReplyFilter(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.Segment{message.Reply(12345), message.At(12345), message.Text("1234")}, NewPattern().Reply().Text("1234"), true},
		{[]message.Segment{message.Reply(12345), message.At(12345), message.Text("1234")}, NewPattern().Reply(true).Text("1234"), false},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
			out := rule(ctx)
			assert.Equal(t, v.expected, out)
		})
	}
}
func TestPattern_Any(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.Segment{message.Text("haha")}, NewPattern().Any(), true},
		{[]message.Segment{message.Image("not a image")}, NewPattern().Any(), true},
		{[]message.Segment{message.At(1919810), message.Reply(12345)}, NewPattern().Any().Reply(), true},
		{[]message.Segment{message.Reply(12345), message.At(1919810)}, NewPattern().Any().At(), true},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
			out := rule(ctx)
			assert.Equal(t, out, v.expected)
		})
	}
	t.Run("get", func(t *testing.T) {
		ctx := fakeCtx([]message.Segment{message.Reply("just for test")})
		rule := NewPattern().Any().AsRule()
		_ = rule(ctx)
		model := PatternModel{}
		err := ctx.Parse(&model)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "just for test", model.Matched[0].Reply())
	})
}
func TestPatternParsed_Gets(t *testing.T) {
	assert.Equal(t, []string{"gaga"}, PatternParsed{value: []string{"gaga"}}.Text())
	assert.Equal(t, "image", PatternParsed{value: "image"}.Image())
	assert.Equal(t, "reply", PatternParsed{value: "reply"}.Reply())
	assert.Equal(t, "114514", PatternParsed{value: "114514"}.At())
	text := message.Text("1234")
	assert.Equal(t, &text, PatternParsed{msg: &text}.Raw())
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
		{[]message.Segment{message.Text("/do it")}, NewPattern().Text("/(do) (.*)").At().SetOptional(true), []PatternParsed{
			{
				value: []string{"/do it", "do", "it"},
			}, {
				value: nil,
			},
		}},
		{[]message.Segment{message.Text("/do it")}, NewPattern().Text("/(do) (.*)").At().SetOptional(false), []PatternParsed{}},
		{[]message.Segment{message.Text("happy bear"), message.At(114514)}, NewPattern().Reply().SetOptional().Text(".+").SetOptional().At().SetOptional(false), []PatternParsed{
			{
				value: nil,
			},
			{
				value: "happy bear",
			},
			{
				value: "114514",
			},
		}},
		{[]message.Segment{message.Text("happy bear"), message.At(114514)}, NewPattern().Image().SetOptional().Image().SetOptional().Image().SetOptional(), []PatternParsed{ // why you do this
			{
				value: nil,
			},
			{
				value: nil,
			},
			{
				value: nil,
			},
		}},
	}
	for i, v := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
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
				assert.Equal(t, v.expected[i].value != nil, parsed.Matched[i].value != nil)
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
		{[]message.Segment{message.Text("test haha test"), message.At(123)}, NewPattern().Text("((ha)+)").At(), []PatternParsed{
			{
				value: []string{"haha", "haha", "ha"},
			}, {
				value: "123",
			},
		}},
		{[]message.Segment{message.Text("haha")}, NewPattern().Text("(h)(a)(h)(a)"), []PatternParsed{
			{
				value: []string{"haha", "h", "a", "h", "a"},
			},
		}},
		{[]message.Segment{message.Reply("fake reply"), message.Image("fake image"), message.At(999), message.At(124), message.Text("haha")}, NewPattern().Reply().Image().At().At(message.NewMessageIDFromInteger(124)).Text("(h)(a)(h)(a)"), []PatternParsed{

			{
				value: "fake reply",
			},
			{
				value: "fake image",
			},
			{
				value: "999",
			},
			{
				value: "124",
			},
			{
				value: []string{"haha", "h", "a", "h", "a"},
			},
		}},
	}
	for i, v := range textTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
			matched := rule(ctx)
			parsed := &PatternModel{}
			err := ctx.Parse(parsed)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, true, matched)
			for i := range parsed.Matched {
				assert.Equal(t, v.expected[i].value, parsed.Matched[i].value)
				assert.Equal(t, &(v.msg[i]), parsed.Matched[i].msg)
			}
		})
	}
}
