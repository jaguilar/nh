package mon

import "github.com/jaguilar/nethack/model/item"

type Monster struct {
	*Species

	// Equipment is the Monster's equipment, with each slot corresponding
	// to the same Bodypart in the Species' anatomy.
	Equipment []*item.Item

	// id is the id of the monster.
	//
	// This engine Calls each monster by a unique id. Monsters
	// with unique names will have those unique names in this field.
	// All others will have the id assigned to them by the engine.
	// We use this to dedup
	id string
}
