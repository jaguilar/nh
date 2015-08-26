package item

import (
	"fmt"
	"strings"
)

// Item is an item in Nethack.
type Item struct {
	BUC
	Enhancement
	Erosion

	// Has the item been Fixed? (Is it rustproof, fireproof, etc.?)
	Fixed bool

	Category

	// Nil if we don't know the exact class (e.g. this is an unidentified potion).
	*Class

	// Called is what we've called the item, if we've called the item's category something.
	Called string

	// Named is what we've named this individual stack of items, if we've named this individual stack
	// something.
	Named string

	Charge

	// The price of an item, if we have price-id'd it.
	// Zeroed for items of a known type, or if we have not yet price id'd the item.
	Price int
}

// BUC is the blessed, cursed, or uncursed status of an item.
type BUC int

// The various BUC statuses. Blessed, Cursed, and Uncursed are all in-game
// statuses. BUCUnknown means we haven't determined the BUC yet. Noncursed
// means we know that it's not cursed (we saw our pet walk on it, maybe).
const (
	BUCUnknown BUC = iota
	Noncursed
	Uncursed
	Blessed
	Cursed
)

// Category is the category of an item.
type Category rune

// The various item categories.
const (
	Amulet        Category = '"'
	Armor                  = '['
	Boulder                = '`'
	HeavyIronBall          = '0'
	Coins                  = '$'
	Comestible             = '%'
	Gem                    = '*'
	IronChain              = '_'
	Potion                 = '!'
	Scroll                 = '?'
	Spellbook              = '+'
	Statue                 = '`'
	Ring                   = '='
	Tool                   = '('
	Wand                   = '/'
	Weapon                 = ')'
)

// MaterialType is the matter that makes up an item.
type MaterialType int

// Various materials. Might not include all materials in the game.
const (
	MaterialUnspecified MaterialType = iota
	Iron
	Copper
	Wood
	Leather
	Cloth
	Plastic
)

// Class is the exact class of an item. For example "silver arrow" is a class.
//
// The model will report an item as having a given class when the class has been
// conclusively inferred, whether or not Nethack knows that the item is of that class.
// It is an invariant that the model should always report the same class for an item
// if the item is subsequently identified and parsed.
//
// An empty string means the exact class is not known.
type Class struct {
	Name string
	MaterialType
}

// Erosion describes the erodedness of an item.
type Erosion struct {
	Rusted, Corroded, Burnt, Rotted ErosionLevel
}

func (e Erosion) String() string {
	var parts []string
	withPrefix := func(e ErosionLevel, t string) string {
		p := e.Prefix()
		if p != "" {
			return p + " " + t
		}
		return t
	}

	// Corrosion needs to be displayed before rust. Other than iron, items cannot
	// have multiple erosions so there is no need to control for the display ordering.
	if e.Rusted != Uneroded {
		parts = append(parts, withPrefix(e.Rusted, "rusty"))
	}
	if e.Corroded != Uneroded {
		parts = append(parts, withPrefix(e.Corroded, "corroded"))
	}
	if e.Burnt != Uneroded {
		parts = append(parts, withPrefix(e.Burnt, "burnt"))
	}
	if e.Rotted != Uneroded {
		parts = append(parts, withPrefix(e.Rotted, "rotted"))
	}

	if parts != nil {
		return strings.Join(parts, " ")
	}
	return ""
}

// ErosionLevel is the extent of a type of erosion. Lower is more eroded.
type ErosionLevel int

// The various levels of erosion.
const (
	ThoroughlyEroded ErosionLevel = -3 + iota
	VeryEroded
	Eroded
	Uneroded
)

// Prefix returns the string prefix conferred by a given ErosionLevel.
// For example, when formatting Rusted at the ThoroughlyEroded level, we want
// to emit "thoroughly rusty", which means this returns "thoroughly".
func (e ErosionLevel) Prefix() string {
	switch e {
	case ThoroughlyEroded:
		return "thoroughly"
	case VeryEroded:
		return "very"
	default:
		return ""
	}
}

// Enhancement represents the enhancement level of an item.
type Enhancement int

// Charge tells the charge state of an item. If the item is of a type that is
// not chargeable, or if the charge is not known, Known will be set to false.
type Charge struct {
	Known        bool
	Cur, Max     int
	TimesCharged int
}

func (c Charge) String() string {
	if !c.Known {
		return ""
	}
	if c.Max < c.Cur {
		return "charged"
	}
	return fmt.Sprintf("(%d:%d)", c.Cur, c.Max)
}
