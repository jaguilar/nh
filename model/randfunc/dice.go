package randfunc

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
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

var diceRe = regexp.MustCompile(`(\d*)d(\d+)`)

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

// Dice converts a "dice expression" like "2d4", return an RFunc that models
// those dice. Each expression must match `\d*d\d+`. That is: any number of
// digits, including zero, followed by the character 'd', followed by at least
// one digit. In practice, if you make these numbers too large, you will have
// overflow problems.
//
// If you have more than
func Dice(s string) (RFunc, error) {
	match := diceRe.FindStringSubmatch(s)
	if match == nil {
		return nil, fmt.Errorf("can't parse as dice expression: %s", s)
	}

	d := dice{n: 1}
	var err error
	if match[1] != "" {
		d.n, err = strconv.Atoi(match[1])
	}
	if err != nil {
		panic(err) // Can't happen: the regexp ensures that we only get ints.
	}
	d.sides, err = strconv.Atoi(match[2])
	if err != nil {
		panic(err) // Can't happen: same reason.
	}
	return d, err
}

// DiceMust is like Dice, but panics if there is an error.
func DiceMust(s string) RFunc {
	d, err := Dice(s)
	if err != nil {
		panic(err)
	}
	return d
}

type constant int

func (c constant) Bound() (int, int) { return int(c), int(c) }
func (c constant) Do() int           { return int(c) }

// Constant returns an RFunc that models a constant number.
func Constant(i int) RFunc {
	return constant(i)
}

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
