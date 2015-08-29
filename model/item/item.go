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

	// The shared item class for this object.
	*Class

	// Named is what we've named this individual stack of items, if we've named this individual stack
	// something.
	Named string

	// The letter this item has in our inventory, if any.
	//
	// Note: if the item is being examined in a scenario where it has no permanent mapping,
	// such as when looting or in a bag, it will probably be expressed as an ephemeral map
	// of rune->*Item. Its inventory letter will be zero.
	InventoryLetter rune

	Charge
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

// Erosion describes the erodedness of an item.
type Erosion struct {
	Rusty, Corroded, Burnt, Rotted ErosionLevel
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
	if e.Rusty != Uneroded {
		parts = append(parts, withPrefix(e.Rusty, "rusty"))
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
	Cur, Max     int
	TimesCharged int
}

func (c Charge) String() string {
	if c.Max == 0 && c.Cur == 0 {
		return "unknown"
	}
	if c.Max < c.Cur {
		return "charged"
	}
	return fmt.Sprintf("(%d:%d)", c.Cur, c.Max)
}
