package model

import (
	"regexp"
	"strings"

	"github.com/jaguilar/nh/model/item"
)

// Screen is the fundamental type we parse. It represents a nethack screen.
// We expect the screen in column-major order.
type Screen [][]rune

// MenuFormat is a type of menu on the Screen. Nethack only has a few discrete
// menu types that we actually need to process.
type MenuFormat int

const (
	// MenuNone - no menu on screen.
	MenuNone MenuFormat = iota

	// MenuInv - screen is showing all or part of the inventory.
	MenuInv

	// MenuSpell - screen is showing the known spell list.
	MenuSpell

	// MenuEnhance - screen is showing the enhance list.
	MenuEnhance

	// MenuUnknown - screen is showing a menu we don't know how to handle.
	MenuUnknown
)

/*
nextMenuContext gets the next menu context from the current one.

Let's look at the format of some nethack menus. We'll use the '·' symbol
to indicate whitespace. The zeroth column of this file is also the zeroth
column of the nethack screen I'm trying to show as an example.

Normal gameplay:

MESSAGE_OR_BLANK
dungeon line 1
dungeon line 2
dungeon line 3

Enhance:

·Pick a skill to advance:
·
·Fighting Skills
·a - bare handed combat

Short Inventory:

              ·WEAPONS                 <- category label on top line, right aligned
dungeon line 1·a - a +1 elven dagger
dungeon line 2·a - a +1 short sword
              ^--- Note: one column left overwritten with blanks.

Long inventory is when the inventory is as long or longer than the screen. In this case,
the inventory is aligned as with #enhance, one column outdented from the very left.

All lists that we handle end with "(end)" or "(x of y)", where the latter form
indicates that there were multiple pages.

The menu state machine is very simple:

         |--|                  |--|
         |  |                  |  |
         -> non-menu <----> menu <-

Each menu type can be reached from non-menu, but you cannot transition from one
menu type to another.
*/
func (s Screen) NextMenu(last MenuFormat) MenuFormat {
	top := string(s[0])
	top = strings.TrimRight(top, " ")

	// If the top line is blank, we're not in a menu.
	// As far as I can tell, there are no exceptions. Also it appears that
	// none of the available menu types fill in the very top left cell. If that
	// is filled, we're printing a status message and thus are no longer in a
	// menu. (We might still be reading user input, however that is handled
	// separately from menus.)
	if top == "" || top[0] != ' ' {
		return MenuNone
	}

	// We never transition from one menu context to another. If we're still
	// in a menu, it's the same one we started in.
	if last != MenuNone {
		return last
	}

	// Ok, so it looks like we're trying to transition to a menu!
	if invListRegexp.MatchString(top) {
		return MenuInv
	}
	if strings.Contains(top, "Pick a skill to advance") || strings.Contains(top, "current skills") {
		return MenuEnhance
	}
	if strings.Contains(top, "Choose which spell to cast") {
		return MenuSpell
	}

	return MenuUnknown
}

var (
	// invListRegexp will match the top line of the inventory if it's just the
	// "i" inventory. Otherwise, it may not, depending on how the inventory was
	// accessed. For example, if you do "d*", the top line will say,
	// "drop what?" For the time being, we will treat such cases as mcUnhandled.
	invListRegexp = regexp.MustCompile(`
		(Amulets)|(Weapons)|(Armor)|
		(Potions)|(Scrolls)|(Wands)|
		(Rings)|(Spellbooks)|(Gems)|
		(Tools)|(Boulders/Statues)|
		(Comestibles)|(Iron balls)`)
)

// ParseItems parses a list of items on the screen. It returns items it was
// able to parse successfully separately from those it could not parse.
// In the long run, we should get to the point where we never return any
// errors.
func (s Screen) ParseItems() (items []*item.Item, errs []string) {
	istrings := s.itemStrings()

	for _, s := range istrings {
		i, err := item.Parse(s)

	}

	return
}

/*
func (g *Game) parseEvents() error {
	var buf bytes.Buffer
	for {
		eventopine := string(g.vt.Content[0])
		eventopine = strings.TrimRight(eventopine, " ")
		if eventopine == "" {
			break
		}
		split := strings.Split(eventopine, "--More--")
		buf.WriteString(split[0]) // Everything before the "--More--" is part of the event.
		if len(split) == 1 {      // Nothing more to show this turn.
			break
		}
	}

	// Currently we interpret no events. We just split after the period in each sentence
	// and append that to the events list. Later, we will begin interpreting these events.
	sentences := strings.SplitAfter(buf.String(), ".")
	for i, s := range sentences {
		sentences[i] = strings.TrimLeft(s, " ")
	}
	g.Events.PushFront(TurnEvents{E: sentences})

	return nil
}

func (g *Game) parseMenu() error {
	return nil
}
*/
