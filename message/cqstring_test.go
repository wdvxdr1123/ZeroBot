package message

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestParseMessageFromString2(t *testing.T) {
	var tests = []struct {
		CQString string
		Expected Message
	}{
		{
			``,
			Message{},
		},
		{
			`Gorilla[CQ:text]`,
			Message{Text("Gorilla"), MessageSegment{Type: "text", Data: map[string]string{}}},
		},
		{
			`[CQ:face,id=123][CQ:face,id=1234]  `,
			Message{Face("123"), Face("1234"), Text("  ")},
		},
		{
			`ȐĉņþƦȻƝƃ[CQ:rcnb][CQ:ɌćƞßɌĆnƅŕĉ,ɌcńƁ=ȓČņÞ]`,
			Message{
				Text("ȐĉņþƦȻƝƃ"),
				MessageSegment{Type: "rcnb", Data: map[string]string{}},
				MessageSegment{Type: "ɌćƞßɌĆnƅŕĉ", Data: map[string]string{"ɌcńƁ": "ȓČņÞ"}},
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := ParseMessageFromString2(test.CQString)
			assert.Equal(t, test.Expected, got)
		})
	}
}

const bench = `[CQ:rcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞrcnbCQȓČņÞ,a=b]`

func BenchmarkParseMessageFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMessageFromString(bench)
	}
}

func BenchmarkParseMessageFromString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMessageFromString2(bench)
	}
}
