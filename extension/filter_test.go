package extension

import (
	"testing"

	"github.com/stretchr/testify/assert"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func TestFilter(t *testing.T) {
	event := zero.Event{
		PostType:    "message",
		DetailType:  "group",
		MessageType: "group",
		SubType:     "abc",
	}

	result := Filter(
		Or(
			PostType("notice"),
			PostType("message"),
		),
		SubType("abc"),
	)(&event, nil)

	assert.Equal(t, true, result)
}
