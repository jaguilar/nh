package item

import (
	"testing"

	"github.com/jaguilar/testify/assert"
)

func TestItemStringer(t *testing.T) {
	// The input is exactly how nethack formats this item.
	i, err := Parse("A - an uncursed very burnt rotted +0 elven cloak")
	if assert.Nil(t, err) {
		// Note: article dropped intentionally.
		assert.Equal(t, "A - uncursed very burnt rotted +0 elven cloak", i.String())
	}
}

// TestItemParseStringParse is the big daddy item regression test. We've selected
// a large number of identified items from dumplogs at the alt.org nethack server.
// We do not expect for Parse().String() to match the input string. However, it
// should be the case that any time we parse the stringified version of an Item,
// the resulting item is the same.
func TestItemParseRegression(t *testing.T) {
	// TODO(jaguilar)
}
