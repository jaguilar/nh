package item

import "io"

func init() {
	// TODO(jaguilar): fill in the possible appearances, load those into the
	// appearance table.
	var data = `name,price,weight,probabilty,eat,appearance
amulet of change,150,20,130c,Y,
amulet of ESP,150,20,175,Y,
amulet of life saving,150,20,75,,
amulet of magical breathing,150,20,65,Y,
amulet of reflection,150,20,75,,
amulet of restful sleep,150,20,135c,Y,
amulet of strangulation,150,20,135c,Y,
amulet of unchanging,150,20,45,Y,
amulet versus poison,150,20,165,Y,
cheap plastic imitation of the Amulet of Yendor,0,20,0,,Amulet of Yendor
Amulet of Yendor,30000,20,0,,`
	csv := mustMapCsv(data)

	for csv.next() {
		c := &Class{
			Category:   Amulet,
			Name:       csv.get("name"),
			Price:      mustInt(csv.get("price")),
			Weight:     mustInt(csv.get("weight")),
			Edible:     "" != csv.get("eat"),
			Appearance: csv.get("appearance"),
		}
		classes[c.Name] = c
	}
	if csv.err != io.EOF {
		panic(csv.err)
	}
}
