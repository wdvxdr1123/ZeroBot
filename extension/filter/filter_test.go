package filter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func TestFilter(t *testing.T) {
	e := map[string]interface{}{
		"post_type": "notice",
		"user_id":   "notice",
	}
	b, _ := json.Marshal(e)
	rawEvent := gjson.ParseBytes(b)
	event := &zero.Event{
		RawEvent: rawEvent,
	}
	result := New(
		func(ctx *zero.Ctx) gjson.Result {
			return ctx.Event.RawEvent
		},
		NewField("post_type").Any(
			Equal("notice"),
			Not(
				In("message"),
			),
		),
		NewField("user_id").All(
			NotEqual("abs"),
		),
	)(&zero.Ctx{Event: event})
	assert.Equal(t, true, result)
}
