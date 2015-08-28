package randfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var convergeTimes = 1000
var near = 0.1

func convergentAvg(r RFunc) float64 {
	var x float64
	for i := 0; i < convergeTimes; i++ {
		x += float64(r.Do())
	}
	return x / float64(convergeTimes)
}

func TestDiceBounds(t *testing.T) {
	for _, tc := range []struct {
		string
		min, max int
		avg      float64
	}{
		{"4d8", 4, 8 * 4, (4 + 8*4) / 2},
		{"d4", 1, 4, float64(1+4) / 2.0},
	} {
		r := DiceMust(tc.string)
		min, max := r.Bound()
		assert.Equal(t, tc.min, min)
		assert.Equal(t, tc.max, max)
		assert.InEpsilon(t, tc.avg, convergentAvg(r), near)
	}
}

func TestDiceRequiresSides(t *testing.T) {
	_, err := Dice("4d")
	assert.NotNil(t, err)
}

func TestConst(t *testing.T) {
	d := Constant(4)
	min, max := d.Bound()
	assert.Equal(t, 4, min)
	assert.Equal(t, 4, max)
	assert.Equal(t, 4, d.Do())
}

func TestSum(t *testing.T) {
	s := Sum(Constant(2), DiceMust("d5"))
	min, max := s.Bound()
	assert.Equal(t, 3, min)
	assert.Equal(t, 7, max)
	assert.InEpsilon(t, 5, convergentAvg(s), near)
}
