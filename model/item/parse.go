package item

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parse parses an item string and returns an appropriate *Item for that string.
func Parse(s string) (*Item, error) {
	m := matchMap(itemRe, itemRe.FindStringSubmatch(s))
	if m == nil {
		return nil, fmt.Errorf("mismatch vs item regexp: %s", s)
	}

	i := &Item{}
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

	return i, nil
}

// alternates returns a wrapped, uncaptured regexp that matches any one of a
// group of alternate regular expressions.
func alternates(inRe ...string) string {
	var anon []string
	for _, s := range inRe {
		anon = append(anon, "(?:"+s+")")
	}
	return "(?:" + strings.Join(anon, "|") + ")"
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
	slot         = "^(?P<slot>[a-zA-Z]) -"
	ordinal      = " (?P<ordinal>" + alternates("an", "a", "the", "\\d{1,2}") + ")"
	buc          = "(?: (?P<buc>" + alternates("blessed", "uncursed", "cursed") + "))?"
	greased      = "(?: (?P<greased>greased))?"
	erosionLevel = "(?: (?P<erosionlevel>" + alternates("thoroughly", "very") + "))?"

	erosionType = "(?: (?P<erosiontype>" + alternates("rusty", "burnt", "corroded", "rotted") + "))"
	oneErosion  = erosionLevel + "?" + erosionType
	erosion     = "(?P<erosions>(?:" + oneErosion + ")*)"

	fixedness = "(?: (?P<fixedness>" + alternates(`\w*proof`, "fixed") + "))?"
	ench      = `(?: (?P<ench>(?:\+|-)\d{1,3}))?`
	itype     = `(?: (?P<type>.+?))`
	called    = `(?: called (?P<called>.+?))?`
	named     = `(?: named (?P<named>.+?))?`
	charge    = `(?: \((?P<charge>\d{1,3}):(?P<maxcharge>\d{1,3})\))?$`

	itemRe = regexp.MustCompile(
		slot + ordinal + buc + greased + erosion + fixedness + ench + itype + called + named + charge)

	erosionRe = regexp.MustCompile(oneErosion)
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
		l := parseErosionLevel(m["erosionlevel"])
		switch m["erosiontype"] {
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
