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

type testModel struct {
	Hello string  `zero:"hello"`
	Pkg   int32   `zero:"pkg"`
	Help  pt      `zero:"help"`
	Love  float64 `zero:"love"`
}

func BenchmarkState_Parse(b *testing.B) {
	var a = &testModel{}
	for i := 0; i < b.N; i++ {
		_ = testState.Parse(a)
	}
}

func TestState_Parse2(t *testing.T) {
	var a = &testModel{}
	_ = testState.Parse(a)
	assert.Equal(t, 520.1314, a.Love)
}
