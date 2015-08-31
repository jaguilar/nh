package item

import (
	"bufio"
	"os"
	"testing"

	"github.com/jaguilar/testify/assert"
)

func TestItemStringer(t *testing.T) {
	// The input is exactly how nethack formats this item.
	i, err := Parse("A - an uncursed very burnt rotted +0 elven cloak")
	if assert.Nil(t, err) {
		// Note: article dropped intentionally.
		assert.Equal(t, "A - uncursed very burnt rotted +0 elven cloak", i.String(), "%#v", i)
	}
}

// TestItemParseStringParse is the big daddy item regression test. We've selected
// a large number of identified items from dumplogs at the alt.org nethack server.
// We do not expect for Parse().String() to match the input string. However, it
// should be the case that any time we parse the stringified version of an Item,
// the resulting item is the same. Any time that doesn't happen, it's a bug and needs
// to be addressed.
//
// The items in this regression test are by no means unique. Once the bot is able
// to play the game, we will begin assembling items it finds into a go-fuzz corpus.
// That corpus will become an exhaustive test of our parsing capabilities.
func TestItemParseRegression(t *testing.T) {
	t.Skip("Cannot yet pass. Need to have the full item corpus first.")
	assert := assert.New(t)
	f, err := os.Open("_testdata/item_regression_test_data.txt")
	if !assert.Nil(err) {
		return
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		s := scanner.Text()
		i, err := Parse(s)
		if !assert.Nil(err, "parse failed: %s", s) {
			continue
		}
		s2 := i.String()
		i2, err := Parse(s2)

		if !assert.Nil(err, "reparse failed: %s orig: %s", s2, s) {
			continue
		}

		c1, c2 := i.Class, i2.Class
		assert.Equal(c1, c2, "%s", s)
		// Nil out the class pointers to avoid them causing test failures.
		i.Class = nil
		i2.Class = nil
		assert.Equal(i, i2, "%s (%s)", s, s2)
	}
	assert.Nil(scanner.Err())
}
