package model

type screenType int

const (
	// screenList: there is a list being rendered on the screen.
	screenList screenType = iota
	// screenMap: there is neither an event or a
	screenMap
)

func (g *Game) parse() {

}

// parseScreenType returns what type of screen is currently on
// the terminal.
//
// Depending on what's in the top line, we can tell what we're
// looking at. If there's a character in the first column,
// we're looking at events. If there's a character in any other
// column, we're parsing some kind of list aligned to that column.
// Otherwise, we're parsing the basic map.
func parseScreenType(screen [][]rune) screenType {
	for i, r := range screen[0] {
		switch {
		case i == 0 && r != ' ':
			return screenEvents
		case r != ' ':
			return screenList
		}
	}
	return screenMap
}
