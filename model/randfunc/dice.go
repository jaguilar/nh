package randfunc

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// RFunc represents a general random calculation from the game. You can get the
// min, max, and average value from it, or you can execute it to get one random
// sample. All nethack random functions have an unskewed probability
// distribution (we think?), so the average can be found by summing min and max
// and dividing by two.
type RFunc interface {
	Bound() (int, int) // Return the inclusive lower and upper bounds of the function.
	Do() int           // Do executes the random function once and returns the result.
}

var diceRe = regexp.MustCompile(`^(\d*?)(d)?(\d+)$`)

// Dice converts a "dice expression" like "2d4 + 1", return an RFunc that models
// those dice.
//
// The grammar is:
//
// DiceExpression       = Expression | Expression "+" DiceExpression .
// Expression           = Dice | Constant .
// Dice                 = [ int_lit ] "d" int_lit
// Constant             = int_lit
//
// Any amount of white space between each expression is ignored. Whitespace
// inside an expression is not allowed. For example:
//
// - "1d2   +1" -- ok
// - "1d2+1"    -- ok
// - "1 d2+1"   -- error
//
// Dice "xdy" means "roll a y sided die x times"
func Dice(s string) (RFunc, error) {
	parts := strings.Split(s, "+")
	var expressions []RFunc

	for _, p := range parts {
		match := diceRe.FindStringSubmatch(p)
		if match == nil {
			return nil, fmt.Errorf("can't parse as dice expression: %s", s)
		}

		isDice := match[2] != ""

		if isDice {
			if match[3] == "" {
				return nil, fmt.Errorf("each dice must have a sides count: %s", s)
			}

			d := dice{n: 1}
			var err error
			if match[1] != "" {
				d.n, err = strconv.Atoi(match[1])
				if err != nil {
					return nil, err
				}
			}
			d.sides, err = strconv.Atoi(match[3])
			if err != nil {
				return nil, err
			}
			expressions = append(expressions, RFunc(d))
		} else { // Constant expression.
			c, err := strconv.Atoi(match[3])
			if err != nil {
				return nil, err
			}
			expressions = append(expressions, RFunc(constant(c)))
		}
	}

	return joined(expressions), nil
}

// DiceMust is like Dice, but panics if there is an error.
func DiceMust(s string) RFunc {
	d, err := Dice(s)
	if err != nil {
		panic(err)
	}
	return d
}

type dice struct {
	n, sides int
}

func (d dice) Bound() (int, int) {
	return d.n * 1, d.n * d.sides
}
func (d dice) Do() int {
	var o int
	for i := 0; i < d.n; i++ {
		o += 1 + rand.Intn(d.sides) // Note: Intn non-inclusive on upper side.
	}
	return o
}

type constant int

func (c constant) Bound() (int, int) { return int(c), int(c) }
func (c constant) Do() int           { return int(c) }

type joined []RFunc

func (j joined) Bound() (int, int) {
	var min, max int
	for _, r := range j {
		rmin, rmax := r.Bound()
		min += rmin
		max += rmax
	}
	return min, max
}

func (j joined) Do() int {
	var d int
	for _, r := range j {
		d += r.Do()
	}
	return d
}

// Sum returns an RFunc that models the sum of other RFuncs.
func Sum(rr ...RFunc) RFunc {
	return joined(rr)
}
