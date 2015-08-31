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

	// Fixed - has the item been fixed (rust/burn/corrode/rot-proofed?).
	Fixed bool

	// Greased - is the item greased?
	Greased bool

	// Class - item class for this object.
	Class *Class

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

	// The size of the stack of this item we have.
	Stack int

	// TODO(jaguilar): Figure out how to represent the candelabrum of invocation.
	// Do we need a map of "special" properties? What about things like comestible
	// age? Maybe an empty interface to support that kind of thing?
}

// String is part of the Stringer interface.
//
// At the moment, we're not going to attempt to reach parity with Nethack's
// printing of item names. Particularly, we're not bothering with the initial
// article.
func (i *Item) String() string {
	var parts []string
	writef := func(format string, args ...interface{}) {
		s := fmt.Sprintf(format, args...)
		if s != "" {
			parts = append(parts, s)
		}
	}

	if i.InventoryLetter != 0 {
		writef("%c -", i.InventoryLetter)
	}

	if i.Stack > 1 {
		writef("%d", i.Stack)
	}

	if i.BUC != BUCUnknown {
		// Note: the model will print out its knowledge of the non-cursed state,
		// even if nethack doesn't display that knowledge.
		writef("%s", strings.ToLower(i.BUC.String()))
	}

	if i.Greased {
		writef("greased")
	}

	writef("%s", i.Erosion.String())

	if i.Fixed {
		writef("%s", i.Class.Material.FixedString())
	}

	if i.Enhancement.Known {
		writef("%s", i.Enhancement.String())
	}

	if !i.Artifact() {
		switch {
		case i.Class.Name != "":
			// Case 1: it's identified. We print neither its appearance or its non-unique name.
			writef("%s", i.Class.Name)
		case i.Class.Called != "":
			// Case 2: it's called something. Don't print its appearance.
			writef("%s called %s", strings.ToLower(i.Class.Category.String()), i.Class.Called)
		default:
			// Case 3: if it's unidentified and uncalled, we print the appearance.
			writef("%s %s", i.Class.Appearance, i.Class.Category)
		}

		// Whether identified, called, or uncalled, we print the named claused after.
		if i.Named != "" {
			writef("named %s", i.Class.Name)
		}
	} else {
		// Identified artifacts are simply printed by their name. They can neither
		// be called nor renamed (at least not in a way that's visible in the UI).
		writef("%s", i.Named)
	}

	// Finally, the charge data, if any.
	if i.Charge.Max != 0 {
		writef("%s", i.Charge.String())
	}

	return strings.Join(parts, " ")
}

// Artifact returns whether item is an artifact. That is, are its name and class the same
// as that of an artifact?
func (i *Item) Artifact() bool {
	return artifacts[artifactKey{name: i.Named, class: i.Class.Name}]
}

// BUC is the blessed, cursed, or uncursed status of an item.
// +gen stringer
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

	// Corrosion needs to be displayed before rust. Burnt before rotted.
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
type Enhancement struct {
	Known bool
	Value int
}

func (e Enhancement) String() string {
	if !e.Known {
		return ""
	}
	return fmt.Sprintf("%+d", e.Value)
}

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
