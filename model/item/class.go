package item

import "github.com/jaguilar/nh/model/randfunc"

// Class is the an Item's grouping. "elven dagger" is a class of items, of which
// there may be a +1 cursed elven dagger, a blessed fixed elven dagger, etc.
//
// At the beginning of the game, we will not know many properties of a class.
// We will know about red potions, but we will not know the Name of a red potion:
// is it a potion of sickness or a potion of gain ability? We also won't know its
// effect when quaffed.
type Class struct {
	Appearance string
	Called     string

	Category

	// Name is the identified name of this item class. Blank if not Identified.
	// You can use the Registry to mark the item as identified even if the
	// middleware has not been able to deduce the item name.
	Name string

	// Price is the base price.
	Price int

	// Weight is the weight of the item.
	Weight int

	// GenProb is the probability that the item will be generated randomly.
	GenProb float64

	Material

	// If this is a weapon, these will be set.
	HitBonus           int
	SmallDam, LargeDam randfunc.RFunc

	// TODO(jaguilar): items should have effects listed for apply, zap, hit, quaff, etc.
	// (Maybe!?)
}

// Category is the category of an item.
type Category rune

// The various item categories.
const (
	Amulet        Category = '"'
	Armor         Category = '['
	Boulder       Category = '`'
	HeavyIronBall Category = '0'
	Coins         Category = '$'
	Comestible    Category = '%'
	Gem           Category = '*'
	IronChain     Category = '_'
	Potion        Category = '!'
	Scroll        Category = '?'
	Spellbook     Category = '+'
	Statue        Category = '`'
	Ring          Category = '='
	Tool          Category = '('
	Wand          Category = '/'
	Weapon        Category = ')'
)

// Material is the matter that makes up an item. All items in a given Class are
// made up of the same material.
type Material string

// Various materials. Might not include all materials in the game.
const (
	Iron    Material = "iron"
	Copper  Material = "copper"
	Metal   Material = "metal"
	Wood    Material = "wood"
	Leather Material = "leather"
	Cloth   Material = "cloth"
	Plastic Material = "plastic"
)