package square

import "strings"

type EngravingQuality int

const (
	Unknown EngravingQuality = iota
	Temporary
	SemiPermanent
	Permanent
)

// Engraving is the engraving on a square. There can be only one. It has
// some text, and a quality, which is unknown by default. Depending on how
// we find the engraving, we may know something about its quality. For example,
// if we find it on a grave, it will be permanent. If we read that it's written
// in the dust, we know it's temporary. If we made it ourselves, that's also
// a way for us to know.
type Engraving struct {
	EngravingQuality
	string
}

func (e Engraving) Elbereths() int {
	return strings.Count(e, "elbereth") + strings.Count(e, "Elbereth")
}
