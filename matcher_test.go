package zero

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatcher_Delete(t *testing.T) {
	OnCommand("").Delete()
	assert.Empty(t, matcherList)
}
