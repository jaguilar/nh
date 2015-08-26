package mon

import (
	"github.com/jaguilar/nh/model/anatomy"
)

type Species struct {
	Name, Class string
	*anatomy.Anatomy
}
