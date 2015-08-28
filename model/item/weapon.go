package item

import (
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/jaguilar/nh/model/randfunc"
)

var data = `
# Source: http://www.steelypips.org/nethack/343/tabs-343.html
# name,cost,weight,probability,material,appearance,+hit,smalldam,largedam
# Dagger
orcish dagger,4,10,.12,iron,crude dagger,2,d3,d3
dagger,4,10,.3,iron,,2,d4,d3
silver dagger,40,12,.03,silver,,2,d4,d3
athame,4,10,0,iron,,2,d4,d3
elven dagger,4,10,.1,wood,,2,d5,d3
# `

var (
	// Weapons is a map of all the weapon names to their associated item Class.
	Weapons map[string]*Class
)

func init() {
	Weapons = make(map[string]*Class)

	csvReader := csv.NewReader(strings.NewReader(data))
	csvReader.Comment = '#'

	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, r := range records {
		c := &Class{
			Category:   Weapon,
			Name:       r[0],
			Price:      mustInt(r[1]),
			Weight:     mustInt(r[2]),
			GenProb:    mustFloat(r[3]),
			Material:   Material(r[4]),
			Appearance: r[5],
			HitBonus:   mustInt(r[6]),
			SmallDam:   randfunc.DiceMust(r[7]),
			LargeDam:   randfunc.DiceMust(r[8]),
		}
		Weapons[c.Name] = c
	}
}

func mustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func mustFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}
