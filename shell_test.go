package zero

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	shellTests := [...]struct {
		shell    string
		expected []string
	}{
		{`rm -rf /*`, []string{"rm", "-rf", "/*"}},
		{`echo   "cat cat" -n`, []string{"echo", "cat cat", "-n"}},
		{`shutdown	halt	init`, []string{"shutdown", "halt", "init"}},
		{`test test2`, []string{"test", "test2"}},
	}
	for i, v := range shellTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out := ParseShell(v.shell)
			assert.Equal(t, v.expected, out)
		})
	}
}

func Test_registerFlag(t *testing.T) {
	type args struct {
		RF    bool   `flag:"rf"`
		File  string `flag:"file"`
		Count int    `flag:"count"`
	}
	got := args{}
	expected := args{
		RF:    true,
		File:  "123",
		Count: 10,
	}
	fs := registerFlag(reflect.TypeOf(args{}), reflect.ValueOf(&got))
	err := fs.Parse([]string{"-rf", "-file=123", "-count", "10"})
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
