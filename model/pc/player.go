package pc

import (
	"github.com/jaguilar/nh/model/item"
	"github.com/jaguilar/nh/model/mon"
)

// +gen stringer
type Alignment int

const (
	Neutral Alignment = iota
	Lawful
	Chaotic
)

type BUC int

const (
	Uncursed BUC = iota
	Blessed
	Cursed
)

type Player struct {
	Name string

	Hp, Pow, HpMax, PowMax       int
	Str, Dex, Con, Int, Wis, Cha int
	Alignment

	*mon.Species

	// Equip is the equipment the Player is wearing. There can be
	// one element for each Bodypart in the player's Species.
	Equip []*item.Item

	Gold int
	Pack []*item.Item
}
