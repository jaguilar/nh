package level

type Branch string

const (
	Dungeon    Branch = "d"
	Mines             = "m"
	Sokoban           = "s"
	Ghennom           = "g"
	Ludios            = "l"
	VladsTower        = "v"
	Planes            = "p"
)

type LevelID struct {
	Branch
	Floor int
}
