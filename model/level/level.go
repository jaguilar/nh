package level

import (
	"github.com/jaguilar/nh/model/mon"
	"github.com/jaguilar/nh/model/square"
)

var (
	SuspectedMonsterLimit = 30
)

type Level struct {
	// LevelID is the unique identifier of the level (branch+floor).
	LevelID

	// Map is all the Squares in the level.
	Map [20][78]square.Square

	// SuspectedMonsters are the monsters we suspect are on this level, but don't
	// know for sure. This includes any monster that was seen once and not killed
	// or otherwise observed to leave the game. This list is distinct from the
	// the monsters in the squares.
	SuspectedMonsters []*mon.Monster
}
