package item

import "fmt"

// Item is an item in Nethack.
type Item struct {
	BUC
	Enhancement
	Erosion
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

type BUC int

const (
	BUCUnknown BUC = iota
	Noncursed
	Uncursed
	Blessed
	Cursed
)

// Category is the category of an item.
type Category rune

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
}

type Erosion int

const (
	ThoroughlyEroded Erosion = -3 + iota
	VeryEroded
	Eroded
	Uneroded
	Fixed
)

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
