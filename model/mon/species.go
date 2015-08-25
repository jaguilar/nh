package mon

type Species struct {
	Name, Class string
	*Anatomy
}

var (
	Humanoid = MonPart{Head, Eyes, Neck, Body, Arms, Hand, Hand, Finger, Finger, Feet}
	// There are many more of these, but for now we've only implemented Humanoid.
)
