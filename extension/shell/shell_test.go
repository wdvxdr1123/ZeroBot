package shell

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	var shellTests = []struct {
		shell    string
		expected []string
		err      error
	}{
		{`rm -rf /*`, []string{"rm", "-rf", "/*"}, nil},
		{`echo   "cat cat" -n`, []string{"echo", "cat cat", "-n"}, nil},
		{`shutdown	halt	init`, []string{"shutdown", "halt", "init"}, nil},
		{`echo "echo`, nil, errors.New("invalid command line string")},
	}
	for _, v := range shellTests {
		t.Run(v.shell, func(t *testing.T) {
			out, err := Parse(v.shell)
			assert.Equal(t, v.err, err)
			assert.Equal(t, v.expected, out)
		})
	}
}
