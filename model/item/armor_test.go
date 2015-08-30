package item

import (
	"testing"

	"github.com/jaguilar/nh/model/anatomy"
	"github.com/jaguilar/testify/assert"
)

func TestArmor(t *testing.T) {
	var mc [4]int                            // We should see stuff from each magic cancellation tier.
	slots := make(map[anatomy.BodyPart]bool) // We should see stuff from each slot.

	for _, a := range classes {
		if a.Category != Armor {
			continue
		}

		mc[a.MagicCancellation]++
		slots[a.Slots[0]] = true
	}

	for i, m := range mc {
		assert.NotZero(t, m, i)
	}
	expectedSlots := []anatomy.BodyPart{
		anatomy.Head, anatomy.Arms, anatomy.Feet,
		anatomy.Torso, anatomy.TorsoOver, anatomy.TorsoUnder,
	}
	for _, s := range expectedSlots {
		assert.True(t, slots[s], s)
	}
}

func TestHindersSpellcasting(t *testing.T) {
	assert.True(t, HindersSpellcasting(classes["dwarvish mithril coat"]))
	assert.False(t, HindersSpellcasting(classes["Hawaiian shirt"]))
}
