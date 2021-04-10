package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		pattern string
		route   *Route
		err     error
	}{
		{
			pattern: "aaaa :    ",
			route:   nil,
			err:     InvalidPattern,
		},
	}

	for _, test := range tests {
		route, err := Parse(test.pattern)
		assert.Equal(t, test.route, route)
		assert.Equal(t, test.err, err)
	}
}
