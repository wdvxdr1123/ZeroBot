package zero

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	testCtx.Event = e1
	assert.Equal(t, true, t1(testCtx))
	assert.Equal(t, true, t2(testCtx))
	assert.Equal(t, false, t3(testCtx))
}

type pt struct {
	x int
	y int
}

var testState = State{
	"hello": "world",
	"pkg":   int32(123),
	"help": pt{
		x: 1,
		y: 2,
	},
	"love": 520.1314,
}

var testCtx = &Ctx{State: testState}

type testModel struct {
	Hello string  `zero:"hello"`
	Pkg   int32   `zero:"pkg"`
	Help  pt      `zero:"help"`
	Love  float64 `zero:"love"`
}

func BenchmarkState_Parse(b *testing.B) {
	a := &testModel{}
	for i := 0; i < b.N; i++ {
		_ = testCtx.Parse(a)
	}
}

func TestState_Parse2(t *testing.T) {
	a := &testModel{}
	assert.NoError(t, testCtx.Parse(a))
	assert.Equal(t, 520.1314, a.Love)
}

func TestMatcher_Delete(t *testing.T) {
	OnCommand("").Delete()
	assert.Empty(t, matcherList)
}
