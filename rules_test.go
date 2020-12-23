package zero

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestType(t *testing.T) {
	t1 := Type("notice/aaa/bbb")
	t2 := Type("notice/aaa")
	t3 := Type("aaa/aaa/bbb")
	e1 := &Event{
		PostType:   "notice",
		DetailType: "aaa",
		SubType:    "bbb",
	}
	assert.Equal(t, true, t1(e1, State{}))
	assert.Equal(t, true, t2(e1, State{}))
	assert.Equal(t, false, t3(e1, State{}))
}
