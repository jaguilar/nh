package item

import "io"

var (
	chargeableRings []string
)

func init() {
	data := `name,price,weight,probability,charge
meat ring,5,1,0,
adornment,100,3,1,Y
hunger,100,3,1c,
protection,100,3,1,Y
protection from shape changers,100,3,1,
stealth,100,3,1,
sustain ability,100,3,1,
warning,100,3,1,
aggravate monster,150,3,1c,
cold resistance,150,3,1,
gain constitution,150,3,1,Y
gain strength,150,3,1,Y
increase accuracy,150,3,1,Y
increase damage,150,3,1,Y
invisibility,150,3,1,
poison resistance,150,3,1,
see invisible,150,3,1,
shock resistance,150,3,1,
fire resistance,200,3,1,
free action,200,3,1,
levitation,200,3,1,
regeneration,200,3,1,
searching,200,3,1,
slow digestion,200,3,1,
teleportation,200,3,1c,
conflict,300,3,1,
polymorph,300,3,1c,
polymorph control,300,3,1,
teleport control,300,3,1,`
	csv := mustMapCsv(data)

	for csv.next() {
		c := &Class{
			Category: Ring,
			Name:     csv.get("name"),
			Price:    mustInt(csv.get("price")),
			Weight:   mustInt(csv.get("weight")),
		}
		classes[c.Name] = c
		if csv.get("charge") != "" {
			chargeableRings = append(chargeableRings, c.Name)
		}
	}
	if csv.err != io.EOF {
		panic(csv.err)
	}
}

// IsChargeableRing returns whether a given ring can be charged with a scroll
// of charging to boost its effect.
func IsChargeableRing(c *Class) bool {
	if c.Category != Ring || c.Name == "" {
		return false
	}

	for _, cr := range chargeableRings {
		if c.Name == cr {
			return true
		}
	}
	return false
}
