package item

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeaponSanity(t *testing.T) {
	// Just a quick sanity check for weapons. Make sure that each weapon has
	// the right category. At least some should have a hit bonus. All should
	// have damage funcs.
	var hit, appearance int
	for _, w := range Weapons {
		assert.NotEqual(t, "", w.Name)
		assert.NotEqual(t, 0, w.Price)
		assert.Equal(t, Weapon, w.Category, fmt.Sprint(w))
		if w.HitBonus != 0 {
			hit++
		}
		if w.Appearance != "" {
			appearance++
		}
		assert.NotNil(t, w.SmallDam)
		assert.NotNil(t, w.LargeDam)
	}
	assert.NotEqual(t, 0, hit)
	assert.NotEqual(t, 0, appearance)
}
