package item

import (
	"fmt"
	"regexp"
	"strconv"
)

// Parse parses an item string and returns an appropriate *Item for that string.
func Parse(s string) (*Item, error) {
	m := matchMap(itemRe, itemRe.FindStringSubmatch(s))
	if m == nil {
		return nil, fmt.Errorf("mismatch vs item regexp: %s", s)
	}

	i := &Item{Class: &Class{}}
	if t, ok := m["type"]; ok {
		i.Class.Name = t
	}

	if c, ok := m["called"]; ok {
		i.Class.Called = c
	}

	if n, ok := m["named"]; ok {
		i.Named = n
	}

	if bucStr, ok := m["buc"]; ok {
		i.BUC = parseBUC(bucStr)
	}

	if enhStr, ok := m["enh"]; ok {
		i.Enhancement.Known = true
		if _, err := fmt.Sscanf(enhStr, "%d", &i.Enhancement.Value); err != nil {
			return nil, fmt.Errorf("error parsing enhancement \"%s\": %v", enhStr, err)
		}
	}

	if charge, ok := m["charge"]; ok {
		maxCharge, ok := m["maxcharge"]
		if !ok {
			return nil, fmt.Errorf("maxcharge not set but charge is: %s", s)
		}
		var err error
		i.Charge.Cur, err = strconv.Atoi(charge)
		if err != nil {
			return nil, err
		}
		i.Charge.Max, err = strconv.Atoi(maxCharge)
		if err != nil {
			return nil, err
		}
	}

	if _, ok := m["fixedness"]; ok {
		// The only way this turns up is if the item is indeed fixed.
		i.Fixed = true
	}

	if erosion, ok := m["erosions"]; ok {
		i.Erosion = parseErosion(erosion)
	}

	if invLetter, ok := m["slot"]; ok {
		i.InventoryLetter = []rune(invLetter)[0]
	}

	return i, nil
}

// matchMap returns a map of each subexpression name to its match, if any.
// A nil return value indicates no match.
func matchMap(re *regexp.Regexp, s []string) map[string]string {
	if s == nil {
		return nil
	}
	m := make(map[string]string)
	for i, subexpName := range re.SubexpNames() {
		if subexpName == "" || s[i] == "" {
			// Unnamed subexpressions are not reported.
			continue
		}

		m[subexpName] = s[i]
	}
	return m
}

// Within each of the below expressions, put spaces at the front of a non-capturing group.
// That way, if any subsequent group does not appear, the match won't fail because of the
// trailing space. And there is always a preceding space except for the slot expression.
//
// Any time there is a general category of things to match, we'll match them all
// with one subexpression, then split the match later.
//
// During the course of this regexp, we capture a little more than we need to in
// some places.
var (
	slot = "^(?:(?P<slot>[a-zA-Z]) -)?"
	// Ordinal is optional so as to support re-parsing.
	ordinal = `\W?(?P<ordinal>an|a|the|\d{1,2})?`
	buc     = `\W?(?P<buc>blessed|uncursed|cursed)?`
	greased = `\W?(?P<greased>greased)?`
	erosion = `\W?(?P<erosion>(?P<elevel>thoroughly|very)?\W?(?P<etype>rusty|burnt|corroded|rotted))`

	// Once we have the erosions string, we'll have to search it with the erosion regexp to parse
	// out each individual erosion.
	erosions = `(?P<erosions>(?:` + erosion + `)*)`

	fixedness = `\W?(?P<fixedness>(?:\w*proof)|fixed)?`
	enh       = `\W?(?P<enh>(?:\+|-)\d)?`
	itype     = `\W?(?P<type>[a-zA-Z ]+?)`
	called    = `(?:\W?called (?P<called>.+?))?`
	named     = `(?:\W?named (?P<named>.+?))?`
	// Match charge info.
	charge = `\W?(?:\((?P<charge>\d{1,3}):(?P<maxcharge>\d{1,3})\))?`
	// Match any non-charge property. Currently ignored. This will match things
	// like the candelabrum of invocation's extra information.
	otherParenProperty = `\W?(?P<paren>(?:\W?\([^)]*\)*?))?$`

	/*
	   Note: we don't take whether the item is being worn into account here.
	   That is information about the item's location in the world, which is not part of
	   parsing the item.
	   	worn = `(?: \((?P<worn>` + alternates(
	   		"being worn", "weapon in \\w*", "in quiver",
	   		"alternate weapon; not wielded", "wielded in other \\w*",
	   		"on \\w* \\w*") + `)\))?`
	*/

	itemRe = regexp.MustCompile(
		slot + ordinal + buc + greased + erosions + fixedness + enh + itype + called + named + charge + otherParenProperty)

	erosionRe = regexp.MustCompile(erosion)
)

// ParseBUC parses a BUC string from nethack. This cannot return Noncursed,
// because that is an inferred status, not one that is displayed in-game.
// Any malformed input will yield BUCUnknown, as will the empty string.
func parseBUC(s string) BUC {
	switch s {
	case "uncursed":
		return Uncursed
	case "blessed":
		return Blessed
	case "cursed":
		return Cursed
	default:
		return BUCUnknown
	}
}

func parseErosion(s string) Erosion {
	var e Erosion
	matches := erosionRe.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return e
	}

	for _, match := range matches {
		m := matchMap(erosionRe, match)
		l := parseErosionLevel(m["elevel"])
		switch m["etype"] {
		case "corroded":
			e.Corroded = l
		case "rusty":
			e.Rusty = l
		case "burnt":
			e.Burnt = l
		case "rotted":
			e.Rotted = l
		}
	}

	return e
}

func parseErosionLevel(s string) ErosionLevel {
	switch s {
	case "thoroughly":
		return ThoroughlyEroded
	case "very":
		return VeryEroded
	default:
		return Eroded
	}
}
