package item

import (
	"fmt"
	"testing"
)

func TestRe(t *testing.T) {
	fmt.Println(itemRe.String())

	fmt.Println(matchMap(itemRe, "a - an uncursed greased thoroughly corroded thoroughly rusty rustproof +100 helm of opposite alignment"))
	fmt.Println(matchMap(itemRe, "b - an uncursed greased thoroughly corroded thoroughly rusty rustproof wand of secret door detection (7:104)"))
}
