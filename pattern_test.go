package zero

import (
	"fmt"
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
	ctx := &Ctx{Event: &Event{Message: msg}, State: map[string]interface{}{}, caller: &messageLogger{
		msgid:  message.NewMessageIDFromInteger(12345),
		caller: mockAPICaller{},
	}}
	return ctx
}

// copy from extension.PatternModel
type PatternModel struct {
	Matched []PatternParsed `zero:"pattern_matched"`
}

// Test Match
func TestPattern_Text(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.Segment{message.Text("/haha")}, NewPattern(nil).Text(`/(\S+)`), true},
		{[]message.Segment{message.Text("haha")}, NewPattern(nil).Text("haha"), true},
		{[]message.Segment{message.Text("aaa")}, NewPattern(nil).Text("not match"), false},
		{[]message.Segment{message.Image("not a image")}, NewPattern(nil).Text("not match"), false},
		{[]message.Segment{message.At(114514)}, NewPattern(nil).Text("not match"), false},
		{[]message.Segment{message.Text("你说的对但是ZeroBot-Plugin 是 ZeroBot 的 实用插件合集")}, NewPattern(nil).Text("实用插件合集"), true},
		{[]message.Segment{message.Text("你说的对但是ZeroBot-Plugin 是 ZeroBot 的 实用插件合集")}, NewPattern(nil).Text("nonono"), false},
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
		{[]message.Segment{message.Text("haha")}, NewPattern(nil).Image(), false},
		{[]message.Segment{message.Text("haha"), message.Image("not a image")}, NewPattern(nil).Image().Image(), false},
		{[]message.Segment{message.Text("haha"), message.Image("not a image")}, NewPattern(nil).Text("haha").Image(), true},
		{[]message.Segment{message.Image("not a image")}, NewPattern(nil).Image(), true},
		{[]message.Segment{message.Image("not a image"), message.Image("not a image")}, NewPattern(nil).Image(), false},
		{[]message.Segment{message.Image("not a image"), message.Image("not a image")}, NewPattern(nil).Image().Image(), true},
		{[]message.Segment{message.Image("not a image"), message.Image("not a image")}, NewPattern(nil).Image().Image().Image(), false},
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

func TestPattern_FuzzyAt(t *testing.T) {
	textTests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected bool
	}{
		{[]message.Segment{message.Text("haha @114514")}, NewPattern(&PatternOption{
			CleanRedundantAt: true,
			FuzzyAt:          true,
		}).Text("haha").At(), true},
		{[]message.Segment{message.Text("haha 114514")}, NewPattern(&PatternOption{
			CleanRedundantAt: true,
			FuzzyAt:          true,
		}).Text("haha").At(), false},
		{[]message.Segment{message.Text("haha @你好")}, NewPattern(&PatternOption{
			CleanRedundantAt: true,
			FuzzyAt:          true,
		}).Text("haha").At(), true},
		{[]message.Segment{message.Text("haha @")}, NewPattern(&PatternOption{
			CleanRedundantAt: true,
			FuzzyAt:          true,
		}).Text("haha").At(), true},
		{[]message.Segment{message.Text("haha @ 你说的对")}, NewPattern(&PatternOption{
			CleanRedundantAt: true,
			FuzzyAt:          true,
		}).Text("haha").At().Text("你说的对"), true},
		{[]message.Segment{message.Text("haha @114514 你说的对")}, NewPattern(&PatternOption{
			CleanRedundantAt: true,
			FuzzyAt:          true,
		}).Text("haha").At().Text("你说的对"), true},
	}
	for _, v := range textTests {
		t.Run(v.msg.String(), func(t *testing.T) {
			ctx := fakeCtx(v.msg)
			rule := v.pattern.AsRule()
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
		{[]message.Segment{message.Text("haha")}, NewPattern(nil).At(), false},
		{[]message.Segment{message.Image("not a image")}, NewPattern(nil).At(), false},
		{[]message.Segment{message.At(114514)}, NewPattern(nil).At(), true},
		{[]message.Segment{message.At(114514)}, NewPattern(nil).At(message.NewMessageIDFromString("1919810")), false},
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
		{[]message.Segment{message.Text("haha")}, NewPattern(nil).Reply(), false},
		{[]message.Segment{message.Image("not a image")}, NewPattern(nil).Reply(), false},
		{[]message.Segment{message.At(1919810), message.Reply(12345)}, NewPattern(nil).Reply().At(), false},
		{[]message.Segment{message.Reply(12345), message.At(1919810)}, NewPattern(nil).Reply().At(), true},
		{[]message.Segment{message.Reply(12345)}, NewPattern(nil).Reply(), true},
		{[]message.Segment{message.Reply(12345), message.At(1919810)}, NewPattern(nil).Reply(), false},
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
		{[]message.Segment{message.Reply(12345), message.At(12345), message.Text("1234")}, NewPattern(nil).Reply().Text("1234"), true},
		{[]message.Segment{message.Reply(12345), message.At(12345), message.Text("1234")}, NewPattern(&PatternOption{
			CleanRedundantAt: false,
			FuzzyAt:          false,
		}).Reply().Text("1234"), false},
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
		{[]message.Segment{message.Text("haha")}, NewPattern(nil).Any(), true},
		{[]message.Segment{message.Image("not a image")}, NewPattern(nil).Any(), true},
		{[]message.Segment{message.At(1919810), message.Reply(12345)}, NewPattern(nil).Any().Reply(), true},
		{[]message.Segment{message.Reply(12345), message.At(1919810)}, NewPattern(nil).Any().At(), true},
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
		rule := NewPattern(nil).Any().AsRule()
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
		NewPattern(nil).SetOptional()
	})
	tests := [...]struct {
		msg      message.Message
		pattern  *Pattern
		expected []PatternParsed
	}{
		{[]message.Segment{message.Text("/do it")}, NewPattern(nil).Text("/(do) (.*)").At().SetOptional(true), []PatternParsed{
			{
				value: []string{"/do it", "do", "it"},
			}, {
				value: nil,
			},
		}},
		{[]message.Segment{message.Text("/do it")}, NewPattern(nil).Text("/(do) (.*)").At().SetOptional(false), []PatternParsed{}},
		{[]message.Segment{message.Text("happy bear"), message.At(114514)}, NewPattern(nil).Reply().SetOptional().Text(".+").SetOptional().At().SetOptional(false), []PatternParsed{
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
		{[]message.Segment{message.Text("happy bear"), message.At(114514)}, NewPattern(nil).Image().SetOptional().Image().SetOptional().Image().SetOptional(), []PatternParsed{ // why you do this
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
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					fmt.Println((parsed.Matched[i].value))
					assert.Equal(t, v.expected[i].value != nil, parsed.Matched[i].value != nil)
				})
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
		{[]message.Segment{message.Text("test haha test"), message.At(123)}, NewPattern(nil).Text("((ha)+)").At(), []PatternParsed{
			{
				value: []string{"haha", "haha", "ha"},
			}, {
				value: "123",
			},
		}},
		{[]message.Segment{message.Text("haha")}, NewPattern(nil).Text("(h)(a)(h)(a)"), []PatternParsed{
			{
				value: []string{"haha", "h", "a", "h", "a"},
			},
		}},
		{[]message.Segment{message.Reply("fake reply"), message.Image("fake image"), message.At(999), message.At(124), message.Text("haha")}, NewPattern(nil).Reply().Image().At().At(message.NewMessageIDFromInteger(124)).Text("(h)(a)(h)(a)"), []PatternParsed{

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
