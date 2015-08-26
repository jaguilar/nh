package item

import (
	"errors"
	"regexp"
	"strings"
)

// Parse parses an item string and returns an appropriate *Item for that string.
func Parse(s string) (*Item, error) {
	return nil, errors.New("unimplemented")
}

// alternates returns an uncaptured regexp that matches any one of a group of alternate
// regular expressions.
func alternates(inRe ...string) string {
	var anon []string
	for _, s := range inRe {
		anon = append(anon, "(?:"+s+")")
	}
	return "(?:" + strings.Join(anon, "|") + ")"
}

// matchMap returns a map of each subexpression name to its match, if any.
// A nil return value indicates no match.
func matchMap(re *regexp.Regexp, s string) map[string]string {
	submatches := re.FindStringSubmatch(s)
	if submatches == nil {
		return nil
	}

	m := make(map[string]string)
	for i, subexpName := range re.SubexpNames() {
		if subexpName == "" || submatches[i] == "" {
			// Unnamed subexpressions are not reported.
			continue
		}

		m[subexpName] = submatches[i]
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
	erosionLevel = alternates(" thoroughly", " very") + "?" // May not appear.
	fixedness    = "(?P<fixedness> " + alternates(`\w*proof`, "fixed") + ")?"
	erosionType  = alternates("rusty", "burnt", "corroded", "rotted")
	oneErosion   = erosionLevel + " " + erosionType + fixedness
	erosions     = "(?P<erosions>(?:" + oneErosion + ")*)"                 // Any number of erosions.
	ench         = `(?: (?P<ench>(?:\+|-)\d{1,3}))?`                       // May not appear.
	itype        = `(?: (?P<type>.+?))`                                    // Non-greedy typename capture.
	charge       = `(?: \((?P<charge>\d{1,3}):(?P<maxcharge>\d{1,3})\))?$` // The charge may or may not appear at the end.

	itemRe = regexp.MustCompile(slot + ordinal + buc + greased + erosions + ench + itype + charge)
)

// ParseBUC parses a BUC string from nethack. This cannot return Noncursed,
// because that is an inferred status, not one that is displayed in-game.
// Any malformed input will yield BUCUnknown, as will the empty string.
func ParseBUC(s string) BUC {
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
