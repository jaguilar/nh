package anatomy

// A Bodypart is any bodypart referred to by Nethack. In practice, we're so far only
// implementing those that can hold pieces of equipment.
type Bodypart int

const (
	Head Bodypart = iota
	Eyes
	Neck
	Body
	Arms
	Hand
	Finger
	Feet
)

// Anatomy is a list of Bodyparts belonging to a Species.
type Anatomy []Bodypart

var (
	Humanoid = MonPart{Head, Eyes, Neck, Body, Arms, Hand, Hand, Finger, Finger, Feet}
	// There are many more of these, but for now we've only implemented Humanoid.
)
