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
			pattern: "test {    ",
			route:   nil,
			err:     InvalidParamName,
		},
		{
			pattern: "123456",
			route: &Route{
				fields: []segment{
					{
						kind:    constPart,
						pattern: "123456",
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		route, err := Parse(test.pattern)
		assert.Equal(t, test.route, route)
		assert.Equal(t, test.err, err)
	}
}

func TestRoute_Match(t *testing.T) {
	var r, _ = Parse(`hello {world}`)
	_, ok := r.Match("hello world")
	assert.True(t, ok)
}
