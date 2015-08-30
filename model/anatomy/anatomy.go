package anatomy

// A BodyPart is any bodypart referred to by Nethack. In practice, we're so far only
// implementing those that can hold pieces of equipment.
type BodyPart string

// A list of various body parts. This is currently not exhaustive.
//
// For the torso parts, we're listing the three separately so we can uniquely
// assign item classes to slots. (Hawaiian shirts go to torso under, etc.)
const (
	Head       BodyPart = "head"
	Eyes       BodyPart = "eyes"
	Neck       BodyPart = "neck"
	TorsoUnder BodyPart = "torso under"
	Torso      BodyPart = "torso"
	TorsoOver  BodyPart = "torso over"
	Arms       BodyPart = "arms" // Arms wear gloves. Don't ask me.
	Hand       BodyPart = "hand"
	Finger     BodyPart = "finger"
	Feet       BodyPart = "feet"
)

// Anatomy is a list of BodyParts belonging to a Species.
type Anatomy []BodyPart

// A set of various anatomies in the game.
var (
	Humanoid = &Anatomy{Head, Eyes, Neck, TorsoUnder, Torso, TorsoOver, Arms, Hand, Hand, Finger, Finger, Feet}
)

var (
	// Blocks describes which BodyParts block other BodyParts when occupied by
	// a cursed item. For example, if you are wearing cursed gloves, fingers
	// are blocked. So Arms blocks Finger. Blocking is transitive.
	Blocks = map[BodyPart]BodyPart{
		Arms:      Finger,
		TorsoOver: Torso,
		Torso:     TorsoUnder,
	}
)
