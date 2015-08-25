package square

import (
	"github.com/jaguilar/nethack/model/item"
	"github.com/jaguilar/nethack/model/mon"
)

// A Square is a single square within a Level. It has a Feature, which is a
// general term that includes both ordinary terrains and dungeon features.
// and can contain a single monster and any number of items.
type Square struct {
	// Feature is the dungeon feature or terrain in this Square.
	// nil should be interpreted as "unexplored".
	Feature

	// The Engraving in this square (the zero type being "no engraving").
	Engraving

	// The Monster in this Square, or nil if no monster.
	//
	// If the player character is in this square and the player character
	// is riding a mount, Monster will be nil, and the mount will be
	// pointed to from the Player object.
	*mon.Monster

	// Items in this square.
	Items []item.Item
}
