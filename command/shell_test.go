package command

import (
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
			out := Parse(v.shell)
			assert.Equal(t, v.expected, out)
		})
	}
}
