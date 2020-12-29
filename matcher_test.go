package zero

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcher_Delete(t *testing.T) {
	OnCommand("").Delete()
	assert.Empty(t, matcherList)
}
