package square

import "fmt"

//+gen stringer
type Terrain int

const (
	Unexplored Terrain = iota
	StaircaseUp
	StaircaseDown
	LadderUp
	LadderDown
	Altar
	Sink
	Fountain
	Tree
	Trap
	DoorClosed
	DoorOpenHoriz
	DoorOpenVert
	DoorOpenDestroyed
	Wall
	Drawbridge
	IronBars
	Floor
	Corridor
	Throne
	Grave
	Water
	Ice
	Lava
	Cloud
	Air
	SolidRock
)

// Feature is a dungeon feature. It includes basic terrain and things
// like sinks, thrones, etc. that need tracking.
type Feature interface {
	fmt.Stringer
}
