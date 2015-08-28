package item

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailedParses(t *testing.T) {
	// TODO(jaguilar): until we have a table of the items in the game, we're
	// not realistically going to be able to fail any tests. This is because the
	// match that captures the item type is very broad.
}

func mustParse(s string) *Item {
	i, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return i
}

func TestParseBuc(t *testing.T) {
	for _, tc := range []struct {
		string
		BUC
	}{
		{"a - an uncursed orcish dagger", Uncursed},
		{"b - a cursed orcish dagger", Cursed},
		{"c - a blessed orcish dagger", Blessed},
		{"d - an orcish dagger", BUCUnknown},
	} {
		assert.Equal(t, tc.BUC, mustParse(tc.string).BUC)
	}
}

func TestParseCharge(t *testing.T) {
	for _, tc := range []struct {
		string
		Charge
	}{
		{"a - a blessed wand of death (1:7)", Charge{Cur: 1, Max: 7}},
		{"b - a blessed wand of death", Charge{}},
	} {
		i, err := Parse(tc.string)
		if assert.Nil(t, err) {
			assert.Equal(t, tc.Charge, i.Charge)
		}
	}
}

func TestParseErosions(t *testing.T) {
	for _, tc := range []struct {
		string
		Erosion
	}{
		{"a - a thoroughly rusty iron helm", Erosion{Rusty: ThoroughlyEroded}},
		{"b - an iron helm", Erosion{}},
		{"c - a very rotted elven cloak", Erosion{Rotted: VeryEroded}},
		{"d - a blessed burnt fireproof +5 pair of water walking boots", Erosion{Burnt: Eroded}},
		{"e - a thoroughly corroded -3 orcish dagger named puddingbane", Erosion{Corroded: ThoroughlyEroded}},
		{"c - a very rotted thoroughly burnt fireproof elven cloak", Erosion{Rotted: VeryEroded, Burnt: ThoroughlyEroded}},
	} {
		i, err := Parse(tc.string)
		if assert.Nil(t, err) {
			assert.Equal(t, tc.Erosion, i.Erosion)
		}
	}
}

func TestParseFixedness(t *testing.T) {
	for _, tc := range []struct {
		string
		bool
	}{
		{"a - a thoroughly rusty iron helm", false},
		{"b - a fixed iron helm", true},
		{"c - a very rotted elven cloak", false},
		{"d - a blessed burnt fireproof +5 pair of water walking boots", true},
		{"e - a thoroughly corroded -3 orcish dagger named puddingbane", false},
		{"c - a very rotted thoroughly burnt fireproof elven cloak", true},
	} {
		i, err := Parse(tc.string)
		if assert.Nil(t, err) {
			assert.Equal(t, tc.bool, i.Fixed)
		}
	}
}

func TestParseSlot(t *testing.T) {
	assert.Equal(t, 'V', mustParse("V - a shortsword").InventoryLetter)
}

func TestNamed(t *testing.T) {
	assert.Equal(t, "shoop da woop", mustParse("d - an elven dagger named shoop da woop").Named)
}

func TestCalled(t *testing.T) {
	assert.Equal(t, "sickness", mustParse("e - a potion called sickness").Class.Called)
}

// TODO(jaguilar): gather a large corpus of items and verify that we can parse them
// correctly. This will be especially important once we are interpreting and limiting the
// item class.
