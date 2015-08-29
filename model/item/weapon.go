package item

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jaguilar/nh/model/randfunc"
)

var data = `# Original source: http://www.steelypips.org/nethack/343/tabs-343.html
# Formatted source: https://docs.google.com/spreadsheets/d/1dbr9QTf0JBFYXpFq-2_ODwLqflQXDy15cPcqBEoT9cw/pubhtml
# Name,Price,Weight,Prob,Pct,Material,Appearance,HitBonus,Small Damage,Large Damage,2h
orcish dagger,4,10,0.12,iron,crude dagger,2,d3,d3,
dagger,4,10,0.3,iron,,2,d4,d3,
silver dagger,40,12,0.03,silver,,2,d4,d3,
athame,4,10,0,iron,,2,d4,d3,
elven dagger,4,10,0.1,wood,runed dagger,2,d5,d3,
worm tooth,2,20,0,,,,d2,d2,
knife (shito),4,5,0.2,iron,,,d3,d2,
stiletto,4,5,0.05,iron,,,d3,d2,
scalpel,6,5,0,metal,,2,d3,d3,
crysknife,100,20,0,mineral,,3,d10,d10,
axe,8,60,0.4,iron,,,d6,d4,
battle-axe,40,120,0.1,iron,double-headed axe,,d8+d4,d6+2d4,1
pick-axe,50,100,0,iron,,,d6,d3,
dwarvish mattock,50,120,0.13,iron,broad pick,-1,d12,d8+2d6,1
orcish short sword,10,30,0.03,iron,crude short sword,,d5,d8,
short sword (wakizashi),10,30,0.08,iron,,,d6,d8,
dwarvish short sword,10,30,0.02,iron,broad short sword,,d7,d8,
elven short sword,10,30,0.02,wood,runed short sword,,d8,d8,
broadsword (ninja-to),10,70,0.08,iron,,,2d4,d6+1,
runesword,300,40,0,iron,runed broadsword,,2d4,d6+1,
elven broadsword,10,70,0.04,wood,runed broadsword,,d6+d4,d6+1,
long sword,15,40,0.5,iron,,,d8,d12,
katana,80,40,0.04,iron,samurai sword,1,d10,d12,
two-handed sword,50,150,0.22,iron,,,d12,3d6,1
tsurugi,500,60,0,metal,long samurai sword,2,d16,d8+2d6,1
scimitar,15,40,0.15,iron,curved sword,,d8,d8,
silver saber,75,40,0.06,silver,,,d8,d8,
club,3,30,0.12,wood,,,d6,d3,
aklys,4,15,0.08,iron,thonged club,,d6,d3,
mace,5,30,0.4,iron,,,d6+1,d6,
morning star,10,120,0.12,iron,,,2d4,d6+1,
flail (nunchaku),4,15,0.4,iron,,,d6+1,2d4,
grappling hook,50,30,0,iron,iron hook,,d2,d6,
war hammer,5,50,0.15,iron,,,d4+1,d4,
quarterstaff,5,40,0.11,wood,staff,,d6,d6,1
partisan,10,80,0.05,iron,vulgar polearm,,d6,d6+1,1
fauchard,5,60,0.06,iron,pole sickle,,d6,d8,1
glaive (naginata),6,75,0.08,iron,single-edged polearm,,d6,d10,1
bec-de-corbin,8,100,0.04,iron,beaked polearm,,d8,d6,1
spetum,5,50,0.05,iron,forked polearm,,d6+1,2d6,1
lucern hammer,7,150,0.05,iron,pronged polearm,,2d4,d6,1
guisarme,5,80,0.06,iron,pruning hook,,2d4,d8,1
ranseur,6,50,0.05,iron,hilted polearm,,2d4,2d4,1
voulge,5,125,0.04,iron,pole cleaver,,2d4,2d4,1
bill-guisarme,7,120,0.04,iron,hooked polearm,,2d4,d10,1
bardiche,7,120,0.04,iron,long poleaxe,,2d4,3d4,1
halberd,10,150,0.08,iron,angled poleaxe,,d10,2d6,1
orcish spear,3,30,0.13,iron,crude spear,,d5,d8,
spear,3,30,0.5,iron,,,d6,d8,
silver spear,40,36,0.02,silver,,,d6,d8,
elven spear,3,30,0.1,wood,runed spear,,d7,d8,
dwarvish spear,3,35,0.12,iron,stout spear,,d8,d8,
javelin,3,20,0.1,iron,throwing spear,,d6,d6,
trident,5,25,0.08,iron,,,d6+1,3d4,
lance,10,180,0.04,iron,,,d6,d8,
orcish bow,60,30,0.12,wood,crude bow,,d2,d2,
bow,60,30,0.24,wood,,,d2,d2,
elven bow,60,30,0.12,wood,runed bow,,d2,d2,
yumi,60,30,0,wood,long bow,,d2,d2,
orcish arrow,2,1,0.2,iron,crude arrow,,d5,d6,
arrow,2,1,0.55,iron,,,d6,d6,
silver arrow,5,1,0.12,silver,,,d6,d6,
elven arrow,2,1,0.2,wood,runed arrow,,d7,d6,
ya,4,1,0.15,metal,bamboo arrow,1,d7,d7,
sling,20,3,0.4,leather,,,d2,d2,
flintstone,1,10,0,mineral,,,d6,d6,
# TODO: handle this. rock|gem|glass,,,0,vary,,,d3,d3,
crossbow,40,50,0.45,wood,,,d2,d2,
crossbow bolt,2,1,0.55,iron,,,d4+1,d6+1,
dart,2,1,0.6,iron,,,d3,d2,
shuriken,5,1,0.35,iron,throwing star,2,d8,d6,
boomerang,20,5,0.15,wood,,,d9,d9,
bullwhip,4,20,0.02,leather,,,d2,1,
rubber hose,3,20,0,PLAS,,,d4,d3,
unicorn horn,100,20,0,BONE,,1,d12,d12,`

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
		name, alt := parseAltName(r[0])

		c := &Class{
			Category:   Weapon,
			Name:       name,
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
		if alt != "" {
			Weapons[alt] = c
		}
	}
}

var altNameRe = regexp.MustCompile(`^(.*?)\W*(\(.*\))?$`)

func parseAltName(orig string) (string, string) {
	matches := altNameRe.FindStringSubmatch(orig)
	if matches == nil {
		panic(fmt.Errorf("failed to find a match for regexp that should match all strings"))
	}
	if matches[1] == "" {
		panic(fmt.Errorf("first submatch should always be non-empty (corpus: %s got: %v)", orig, matches))
	}

	return matches[1], matches[2]
}

func mustInt(s string) int {
	if s == "" {
		return 0
	}
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
