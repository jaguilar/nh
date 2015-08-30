package item

import (
	"io"

	"github.com/jaguilar/nh/model/anatomy"
)

func init() {
	// Currently don't parse eff or appearances preceded by *.
	armorData := `name,price,weight,probability,ac,material,effect,mc,appearance,slot
Hawaiian shirt,3,5,8,0,cloth,Shops,,,torso under
T-shirt,2,5,2,0,cloth,Shops,,,torso under
leather jacket,10,30,12,1,leather,,,,torso
leather armor,5,150,82,2,leather,,,,torso
orcish ring mail,80,250,20,2,iron,#,1,crude ring mail,torso
studded leather armor,15,200,72,3,leather,#,1,,torso
ring mail,100,250,72,3,iron,,,,torso
scale mail,45,250,72,4,iron,,,,torso
orcish chain mail,75,300,20,4,iron,#,1,crude chain mail,torso
chain mail,75,300,72,5,iron,#,1,,torso
elven mithril coat,240,150,15,5,mithril,###,3,,torso
splint mail,80,400,62,6,iron,#,1,,torso
banded mail,90,350,72,6,iron,,,,torso
dwarvish mithril coat,240,150,10,6,mithril,###,3,,torso
bronze plate mail,400,450,25,6,copper,,,,torso
plate mail (tanko),600,450,44,7,iron,##,2,,torso
crystal plate mail,820,450,10,7,glass,##,2,,torso
red dragon scales,500,40,0,3,dragon,Fire,,,torso
white dragon scales,500,40,0,3,dragon,Cold,,,torso
orange dragon scales,500,40,0,3,dragon,Sleep,,,torso
blue dragon scales,500,40,0,3,dragon,Elec,,,torso
green dragon scales,500,40,0,3,dragon,Poison,,,torso
yellow dragon scales,500,40,0,3,dragon,Acd,,,torso
black dragon scales,700,40,0,3,dragon,Disint,,,torso
silver dragon scales,700,40,0,3,dragon,Reflect,,,torso
gray dragon scales,700,40,0,3,dragon,Magic,,,torso
red dragon scale mail,900,40,0,9,dragon,Fire,,,torso
white dragon scale mail,900,40,0,9,dragon,Cold,,,torso
orange dragon scale mail:,900,40,0,9,dragon,Sleep,,,torso
blue dragon scale mail,900,40,0,9,dragon,Elec,,,torso
green dragon scale mail,900,40,0,9,dragon,Poison,,,torso
yellow dragon scale mail,900,40,0,9,dragon,Acd,,,torso
black dragon scale mail,1200,40,0,9,dragon,Disint,,,torso
silver dragon scale mail,1200,40,0,9,dragon,Reflect,,,torso
gray dragon scale mail,1200,40,0,9,dragon,Magic,,,torso
mummy wrapping,2,3,0,0,cloth,#Vis,1,,torso over
orcish cloak,40,10,8,0,cloth,##,2,coarse mantelet,torso over
dwarvish cloak,50,10,8,0,cloth,##,2,hooded cloak,torso over
leather cloak,40,15,8,1,leather,#,1,,torso over
cloak of displacement,50,10,10,1,cloth,##Displ,2,*piece of cloth,torso over
oilskin cloak,50,10,10,1,cloth,###Water,3,slippery cloak,torso over
alchemy smock,50,10,9,1,cloth,#Poi+Acd,1,apron,torso over
cloak of invisibility,60,10,10,1,cloth,##Invis,2,*opera cloak,torso over
clk of magic resistance,60,10,2,1,cloth,###Magic,3,*ornamental cope,torso over
elven cloak,60,10,8,1,cloth,###Stlth,3,faded pall,torso over
robe,50,15,3,2,cloth,###Spell,3,,torso over
cloak of protection,50,10,9,3,cloth,###Prot,3,*tattered cape,torso over
fedora,1,3,0,0,cloth,,,,head
dunce cap,1,4,3,0,cloth,Stupid,,conical hat,head
cornuthaum,80,4,3,0,cloth,##Clair,2,conical hat,head
dented pot,8,10,2,1,iron,,,,head
elven leather helm,8,3,6,1,leather,,,leather hat,head
helmet (kabuto),10,30,10,1,iron,,,*plumed helmet,head
orcish helm,10,30,6,1,iron,,,iron skull cap,head
helm of brilliance,50,50,6,1,iron,Int+Wis,,*etched helmet,head
helm of opposite alignment,50,50,6c,1,iron,Align,,*crested helmet,head
helm of telepathy,50,50,2,1,iron,ESP,,*visored helmet,head
dwarvish iron helm,20,40,6,2,iron,,,hard hat,head
leather gloves (yugake),8,10,16,1,leather,,,*old gloves,arms
gauntlets of dexterity,50,10,8,1,leather,Dex,,*padded gloves,arms
gauntlets of fumbling,50,10,8c,1,leather,Fumble,,*riding gloves,arms
gauntlets of power,50,30,8,1,iron,Str,,*fencing gloves,arms
small shield,3,30,6,1,wood,,,,hand
orcish shield,7,50,2,1,iron,,,red-eyed,hand
Uruk-hai shield,7,50,2,1,iron,,,white-handed,hand
elven shield,7,40,2,2,wood,,,blue and green,hand
dwarvish roundshield,10,100,4,2,iron,,,large round,hand
large shield,10,100,7,2,iron,,,,hand
shield of reflection,50,50,3,2,silver,Reflect,,polished silver shield,hand
low boots,8,10,25,1,leather,,,walking shoes,feet
elven boots,8,15,12,1,leather,Stlth,,*mud boots,feet
kicking boots,8,15,12,1,iron,Kick,,*buckled boots,feet
fumble boots,30,20,12c,1,leather,Fumble,,*riding boots,feet
levitation boots,30,15,12c,1,leather,Lev,,*snow boots,feet
jumping boots,50,20,12,1,leather,Jump,,*hiking boots,feet
speed boots,50,20,12,1,leather,Speed,,*combat boots,feet
water walking boots,50,20,12,1,leather,WWalk,,*jungle boots,feet
high boots,12,20,15,2,leather,,,jackboots,feet
iron shoes,16,50,7,2,iron,,,hard shoes,feet`
	csv := mustMapCsv(armorData)

	for csv.next() {
		name, alt := parseAltName(csv.get("name"))
		appearance := csv.get("appearance")
		if appearance != "" && appearance[0] == '*' {
			// * in appearance in this table indicates that it is one of several
			// alternates. We will eventually load these into a list of alternates,
			// but for now we will simply ignore them.
			appearance = ""
		}

		c := &Class{
			Category:          Armor,
			Name:              name,
			Price:             mustInt(csv.get("price")),
			Weight:            mustInt(csv.get("weight")),
			AC:                mustInt(csv.get("ac")),
			Material:          Material(csv.get("material")),
			MagicCancellation: mustInt(csv.get("mc")),
			Slots:             []anatomy.BodyPart{anatomy.BodyPart(csv.get("slot"))},
		}
		classes[c.Name] = c
		classes[alt] = c
	}
	if csv.err != io.EOF {
		panic(csv.err)
	}
}
